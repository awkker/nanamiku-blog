package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultGitHubCacheTTL = 30 * time.Minute

type GitHubProfileConfig struct {
	APIToken   string
	APIBaseURL string
	CacheTTL   time.Duration
}

type GitHubProfileService struct {
	rdb        *redis.Client
	apiToken   string
	apiBaseURL string
	cacheTTL   time.Duration
	client     *http.Client
}

func NewGitHubProfileService(rdb *redis.Client, cfg GitHubProfileConfig) *GitHubProfileService {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.APIBaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}

	cacheTTL := cfg.CacheTTL
	if cacheTTL <= 0 {
		cacheTTL = defaultGitHubCacheTTL
	}

	return &GitHubProfileService{
		rdb:        rdb,
		apiToken:   strings.TrimSpace(cfg.APIToken),
		apiBaseURL: baseURL,
		cacheTTL:   cacheTTL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

type GitHubProfilePayload struct {
	Profile      GitHubProfileSummary  `json:"profile"`
	TechStack    []GitHubTechStackItem `json:"tech_stack"`
	RecentRepos  []GitHubRecentRepo    `json:"recent_repos"`
	ActivityData []int                 `json:"activity_data"`
}

type GitHubProfileSummary struct {
	AvatarURL  string `json:"avatar_url"`
	Name       string `json:"name"`
	Bio        string `json:"bio"`
	HTMLURL    string `json:"html_url"`
	TotalRepos int    `json:"total_repos"`
	TotalStars int    `json:"total_stars"`
	Followers  int    `json:"followers"`
}

type GitHubTechStackItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type GitHubRecentRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stars       int    `json:"stars"`
	PushedAt    string `json:"pushed_at"`
}

type gitHubUser struct {
	AvatarURL   string `json:"avatar_url"`
	Name        string `json:"name"`
	Bio         string `json:"bio"`
	HTMLURL     string `json:"html_url"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
}

type gitHubRepo struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	HTMLURL         string `json:"html_url"`
	Language        string `json:"language"`
	StargazersCount int    `json:"stargazers_count"`
	PushedAt        string `json:"pushed_at"`
}

type gitHubEvent struct {
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	Payload   struct {
		Size    int             `json:"size"`
		Commits json.RawMessage `json:"commits"`
	} `json:"payload"`
}

func (s *GitHubProfileService) GetProfile(ctx context.Context, username string) (*GitHubProfilePayload, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("github username required")
	}

	if cached, ok := s.readCache(ctx, username); ok {
		return cached, nil
	}

	var user gitHubUser
	if err := s.fetchJSON(ctx, fmt.Sprintf("/users/%s", url.PathEscape(username)), &user); err != nil {
		return nil, fmt.Errorf("fetch github user: %w", err)
	}

	repos, err := s.fetchAllPublicRepos(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("fetch github repos: %w", err)
	}

	activityData := make([]int, 12)
	if points, err := s.fetchCommitActivityByMonth(ctx, username); err == nil {
		activityData = points
	}

	totalStars := 0
	langCount := make(map[string]int)
	for _, repo := range repos {
		totalStars += repo.StargazersCount
		if lang := strings.TrimSpace(repo.Language); lang != "" {
			langCount[lang]++
		}
	}

	sort.Slice(repos, func(i, j int) bool {
		return parseTime(repos[i].PushedAt).After(parseTime(repos[j].PushedAt))
	})

	techStack := make([]GitHubTechStackItem, 0, len(langCount))
	for name, count := range langCount {
		techStack = append(techStack, GitHubTechStackItem{Name: name, Count: count})
	}
	sort.Slice(techStack, func(i, j int) bool {
		if techStack[i].Count == techStack[j].Count {
			return techStack[i].Name < techStack[j].Name
		}
		return techStack[i].Count > techStack[j].Count
	})
	if len(techStack) > 10 {
		techStack = techStack[:10]
	}

	recentRepos := make([]GitHubRecentRepo, 0, minInt(len(repos), 6))
	for _, repo := range repos[:minInt(len(repos), 6)] {
		recentRepos = append(recentRepos, GitHubRecentRepo{
			Name:        repo.Name,
			Description: repo.Description,
			HTMLURL:     repo.HTMLURL,
			Language:    repo.Language,
			Stars:       repo.StargazersCount,
			PushedAt:    repo.PushedAt,
		})
	}

	payload := &GitHubProfilePayload{
		Profile: GitHubProfileSummary{
			AvatarURL:  user.AvatarURL,
			Name:       firstNonEmpty(user.Name, username),
			Bio:        user.Bio,
			HTMLURL:    firstNonEmpty(user.HTMLURL, "https://github.com/"+username),
			TotalRepos: maxInt(user.PublicRepos, len(repos)),
			TotalStars: totalStars,
			Followers:  user.Followers,
		},
		TechStack:    techStack,
		RecentRepos:  recentRepos,
		ActivityData: activityData,
	}

	s.writeCache(ctx, username, payload)
	return payload, nil
}

func (s *GitHubProfileService) fetchAllPublicRepos(ctx context.Context, username string) ([]gitHubRepo, error) {
	repos := make([]gitHubRepo, 0, 64)
	const perPage = 100

	for page := 1; page <= 10; page++ {
		var batch []gitHubRepo
		path := fmt.Sprintf("/users/%s/repos?sort=pushed&per_page=%d&page=%d", url.PathEscape(username), perPage, page)
		if err := s.fetchJSON(ctx, path, &batch); err != nil {
			return nil, err
		}
		repos = append(repos, batch...)
		if len(batch) < perPage {
			break
		}
	}

	return repos, nil
}

func (s *GitHubProfileService) fetchCommitActivityByMonth(ctx context.Context, username string) ([]int, error) {
	counts := make([]int, 12)
	now := time.Now().UTC()
	const perPage = 100

	for page := 1; page <= 3; page++ {
		var events []gitHubEvent
		path := fmt.Sprintf("/users/%s/events/public?per_page=%d&page=%d", url.PathEscape(username), perPage, page)
		if err := s.fetchJSON(ctx, path, &events); err != nil {
			return nil, err
		}
		if len(events) == 0 {
			break
		}

		for _, event := range events {
			if event.Type != "PushEvent" {
				continue
			}
			ts := parseTime(event.CreatedAt)
			if ts.IsZero() {
				continue
			}
			diff := (now.Year()-ts.Year())*12 + int(now.Month()-ts.Month())
			if diff < 0 || diff >= 12 {
				continue
			}

			increment := event.Payload.Size
			if increment <= 0 {
				var commits []json.RawMessage
				_ = json.Unmarshal(event.Payload.Commits, &commits)
				if len(commits) > 0 {
					increment = len(commits)
				} else {
					increment = 1
				}
			}

			counts[11-diff] += increment
		}

		if len(events) < perPage {
			break
		}
	}

	return counts, nil
}

func (s *GitHubProfileService) fetchJSON(ctx context.Context, path string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.apiBaseURL+path, nil)
	if err != nil {
		return fmt.Errorf("create github request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "nanamiku-blog/1.0")
	if s.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiToken)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("request github api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("github api returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode github response: %w", err)
	}
	return nil
}

func (s *GitHubProfileService) readCache(ctx context.Context, username string) (*GitHubProfilePayload, bool) {
	if s.rdb == nil {
		return nil, false
	}

	cached, err := s.rdb.Get(ctx, githubProfileCacheKey(username)).Bytes()
	if err != nil {
		return nil, false
	}

	var payload GitHubProfilePayload
	if err := json.Unmarshal(cached, &payload); err != nil {
		return nil, false
	}
	return &payload, true
}

func (s *GitHubProfileService) writeCache(ctx context.Context, username string, payload *GitHubProfilePayload) {
	if s.rdb == nil || payload == nil {
		return
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_ = s.rdb.Set(ctx, githubProfileCacheKey(username), encoded, s.cacheTTL).Err()
}

func githubProfileCacheKey(username string) string {
	return "github:profile:" + strings.ToLower(strings.TrimSpace(username))
}

func parseTime(raw string) time.Time {
	ts, _ := time.Parse(time.RFC3339, raw)
	return ts
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
