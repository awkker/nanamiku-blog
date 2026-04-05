package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const defaultSmokeBaseURL = "http://127.0.0.1:8080/api/v1"

type smokeConfig struct {
	BaseURL         string
	AdminIdentifier string
	AdminPassword   string
	Timeout         time.Duration
}

type apiClient struct {
	baseURL string
	client  *http.Client
}

type apiEnvelope[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type pagedData[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

type loginResponse struct {
	ExpiresAt int64 `json:"expires_at"`
}

type meResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type guestbookMessageResponse struct {
	ID string `json:"id"`
}

type adminGuestbookMessage struct {
	ID         string `json:"id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	Status     string `json:"status"`
}

type friendApplicationSubmission struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type adminFriendApplication struct {
	ID      string `json:"id"`
	SiteURL string `json:"site_url"`
	Status  string `json:"status"`
}

type approveFriendApplicationResponse struct {
	FriendLinkID string `json:"friend_link_id"`
}

type adminFriendLink struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

type dashboardStats struct {
	PendingComments int64 `json:"pending_comments"`
	DraftCount      int64 `json:"draft_count"`
}

type momentItem struct {
	ID           string          `json:"id"`
	Content      string          `json:"content"`
	LikeCount    int64           `json:"like_count"`
	RepostCount  int64           `json:"repost_count"`
	CommentCount int64           `json:"comment_count"`
	Liked        bool            `json:"liked"`
	Reposted     bool            `json:"reposted"`
	Comments     []momentComment `json:"comments"`
}

type momentComment struct {
	ID        string `json:"id"`
	Author    string `json:"author_name"`
	Content   string `json:"content"`
	LikeCount int64  `json:"like_count"`
	Liked     bool   `json:"liked"`
}

func main() {
	_ = godotenv.Load()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	adminAPI, err := newAPIClient(cfg.BaseURL, cfg.Timeout)
	if err != nil {
		log.Fatalf("create admin client: %v", err)
	}
	publicAPI, err := newAPIClient(cfg.BaseURL, cfg.Timeout)
	if err != nil {
		log.Fatalf("create public client: %v", err)
	}

	runID := strings.ToLower(strings.Split(uuid.NewString(), "-")[0])
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := runSmoke(ctx, cfg, runID, adminAPI, publicAPI); err != nil {
		log.Fatalf("smoke failed: %v", err)
	}

	log.Printf("smoke passed: run_id=%s", runID)
}

func loadConfig() (*smokeConfig, error) {
	baseURLFlag := flag.String("base-url", "", "API base URL (default: http://127.0.0.1:8080/api/v1)")
	identifierFlag := flag.String("identifier", "", "admin identifier (username or email)")
	passwordFlag := flag.String("password", "", "admin password")
	timeoutFlag := flag.Duration("timeout", 45*time.Second, "overall smoke timeout")
	flag.Parse()

	cfg := &smokeConfig{
		BaseURL:         normalizeBaseURL(firstNonEmpty(*baseURLFlag, os.Getenv("SMOKE_BASE_URL"), defaultSmokeBaseURL)),
		AdminIdentifier: strings.TrimSpace(firstNonEmpty(*identifierFlag, os.Getenv("SMOKE_ADMIN_IDENTIFIER"))),
		AdminPassword:   strings.TrimSpace(firstNonEmpty(*passwordFlag, os.Getenv("SMOKE_ADMIN_PASSWORD"))),
		Timeout:         *timeoutFlag,
	}

	if cfg.AdminIdentifier == "" || cfg.AdminPassword == "" {
		return nil, fmt.Errorf("missing admin credentials: pass -identifier/-password or set SMOKE_ADMIN_IDENTIFIER and SMOKE_ADMIN_PASSWORD")
	}
	if cfg.Timeout <= 0 {
		return nil, fmt.Errorf("timeout must be greater than 0")
	}
	return cfg, nil
}

func runSmoke(ctx context.Context, cfg *smokeConfig, runID string, adminAPI, publicAPI *apiClient) error {
	log.Printf("smoke: login as %s", cfg.AdminIdentifier)
	loginData, err := doJSON[loginResponse](ctx, adminAPI, http.MethodPost, "/auth/login", map[string]string{
		"identifier": cfg.AdminIdentifier,
		"password":   cfg.AdminPassword,
	})
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}
	log.Printf("smoke: login ok, expires_at=%d", loginData.ExpiresAt)

	me, err := doJSON[meResponse](ctx, adminAPI, http.MethodGet, "/auth/me", nil)
	if err != nil {
		return fmt.Errorf("load current admin: %w", err)
	}
	if strings.TrimSpace(me.Username) == "" {
		return fmt.Errorf("current admin response missing username")
	}
	log.Printf("smoke: /auth/me ok, username=%s", me.Username)

	if _, err := doJSON[dashboardStats](ctx, adminAPI, http.MethodGet, "/admin/dashboard/stats", nil); err != nil {
		return fmt.Errorf("dashboard stats probe: %w", err)
	}
	if _, err := doPagedJSON[map[string]any](ctx, adminAPI, http.MethodGet, "/admin/posts?page=1&size=5", nil); err != nil {
		return fmt.Errorf("admin posts probe: %w", err)
	}
	if _, err := doPagedJSON[map[string]any](ctx, adminAPI, http.MethodGet, "/admin/moments?page=1&size=5", nil); err != nil {
		return fmt.Errorf("admin moments probe: %w", err)
	}
	if _, err := doPagedJSON[map[string]any](ctx, adminAPI, http.MethodGet, "/admin/comments?page=1&size=5", nil); err != nil {
		return fmt.Errorf("admin comments probe: %w", err)
	}

	guestbookID, err := exerciseGuestbookFlow(ctx, runID, adminAPI, publicAPI)
	if err != nil {
		return err
	}
	log.Printf("smoke: guestbook moderation flow ok, message_id=%s", guestbookID)

	friendLinkID, err := exerciseFriendApplicationFlow(ctx, runID, adminAPI, publicAPI)
	if err != nil {
		return err
	}
	log.Printf("smoke: friend application flow ok, friend_link_id=%s", friendLinkID)

	momentID, err := exerciseMomentFlow(ctx, runID, adminAPI, publicAPI)
	if err != nil {
		return err
	}
	log.Printf("smoke: moments interaction flow ok, moment_id=%s", momentID)

	if err := doNoData(ctx, adminAPI, http.MethodPost, "/auth/logout", map[string]string{}); err != nil {
		return fmt.Errorf("logout: %w", err)
	}
	if err := expectStatus(ctx, adminAPI, http.MethodGet, "/auth/me", http.StatusUnauthorized); err != nil {
		return fmt.Errorf("post-logout auth check: %w", err)
	}
	log.Printf("smoke: logout invalidated session")

	return nil
}

func exerciseGuestbookFlow(ctx context.Context, runID string, adminAPI, publicAPI *apiClient) (string, error) {
	nickname := "Smoke Guest " + runID
	content := "Smoke guestbook message " + runID

	created, err := doJSON[guestbookMessageResponse](ctx, publicAPI, http.MethodPost, "/guestbook/messages", map[string]string{
		"author_name":    nickname,
		"author_website": "https://example.com/smoke/" + runID,
		"content":        content,
	})
	if err != nil {
		return "", fmt.Errorf("create guestbook message: %w", err)
	}

	pendingMessages, err := doPagedJSON[adminGuestbookMessage](ctx, adminAPI, http.MethodGet, "/admin/guestbook/messages?status=pending&page=1&size=100", nil)
	if err != nil {
		return "", fmt.Errorf("list pending guestbook messages: %w", err)
	}
	if !containsGuestbookMessage(pendingMessages.Items, created.ID, nickname) {
		return "", fmt.Errorf("pending guestbook queue missing message %s", created.ID)
	}

	if err := doNoData(ctx, adminAPI, http.MethodPost, "/admin/guestbook/messages/"+created.ID+"/approve", map[string]string{}); err != nil {
		return "", fmt.Errorf("approve guestbook message %s: %w", created.ID, err)
	}

	publicMessages, err := doPagedJSON[guestbookMessageResponse](ctx, publicAPI, http.MethodGet, "/guestbook/messages?sort=newest&page=1&size=100", nil)
	if err != nil {
		return "", fmt.Errorf("list public guestbook messages: %w", err)
	}
	if !containsGuestbookID(publicMessages.Items, created.ID) {
		return "", fmt.Errorf("approved guestbook message %s not visible publicly", created.ID)
	}

	if err := doNoData(ctx, adminAPI, http.MethodDelete, "/admin/guestbook/messages/"+created.ID, nil); err != nil {
		return "", fmt.Errorf("cleanup guestbook message %s: %w", created.ID, err)
	}

	return created.ID, nil
}

func exerciseFriendApplicationFlow(ctx context.Context, runID string, adminAPI, publicAPI *apiClient) (string, error) {
	siteURL := "https://smoke-" + runID + ".example.com"
	submission, err := doJSON[friendApplicationSubmission](ctx, publicAPI, http.MethodPost, "/friends/applications", map[string]string{
		"site_name":     "Smoke Friend " + runID,
		"site_url":      siteURL,
		"avatar_url":    "https://example.com/avatar-" + runID + ".png",
		"description":   "Smoke friend application " + runID,
		"contact_email": "smoke+" + runID + "@example.com",
		"contact_note":  "Created by smoke run " + runID,
	})
	if err != nil {
		return "", fmt.Errorf("create friend application: %w", err)
	}
	if submission.Status != "pending" {
		return "", fmt.Errorf("friend application %s returned unexpected status %q", submission.ID, submission.Status)
	}

	applications, err := doPagedJSON[adminFriendApplication](ctx, adminAPI, http.MethodGet, "/admin/friends/applications?status=pending&page=1&size=100", nil)
	if err != nil {
		return "", fmt.Errorf("list pending friend applications: %w", err)
	}
	if !containsFriendApplication(applications.Items, submission.ID, siteURL) {
		return "", fmt.Errorf("pending friend application queue missing application %s", submission.ID)
	}

	approved, err := doJSON[approveFriendApplicationResponse](ctx, adminAPI, http.MethodPost, "/admin/friends/applications/"+submission.ID+"/approve", map[string]any{
		"sort_order":  999,
		"review_note": "smoke approve " + runID,
	})
	if err != nil {
		return "", fmt.Errorf("approve friend application %s: %w", submission.ID, err)
	}
	if strings.TrimSpace(approved.FriendLinkID) == "" {
		return "", fmt.Errorf("approve friend application %s returned empty friend_link_id", submission.ID)
	}

	links, err := doPagedJSON[adminFriendLink](ctx, adminAPI, http.MethodGet, "/admin/friends?page=1&size=100", nil)
	if err != nil {
		return "", fmt.Errorf("list admin friends: %w", err)
	}
	if !containsFriendLink(links.Items, approved.FriendLinkID, siteURL) {
		return "", fmt.Errorf("approved friend link %s not visible in admin list", approved.FriendLinkID)
	}

	if err := doNoData(ctx, adminAPI, http.MethodDelete, "/admin/friends/"+approved.FriendLinkID, nil); err != nil {
		return "", fmt.Errorf("cleanup friend link %s: %w", approved.FriendLinkID, err)
	}

	return approved.FriendLinkID, nil
}

func exerciseMomentFlow(ctx context.Context, runID string, adminAPI, publicAPI *apiClient) (string, error) {
	content := "Smoke moment " + runID
	created, err := doJSON[momentItem](ctx, adminAPI, http.MethodPost, "/admin/moments", map[string]any{
		"content":        content,
		"image_urls":     []string{},
		"publish_status": "published",
	})
	if err != nil {
		return "", fmt.Errorf("create moment: %w", err)
	}
	if strings.TrimSpace(created.ID) == "" {
		return "", fmt.Errorf("create moment returned empty id")
	}

	cleanupErr := func(cause error) error {
		if err := doNoData(ctx, adminAPI, http.MethodDelete, "/admin/moments/"+created.ID, nil); err != nil {
			return fmt.Errorf("%w (cleanup moment %s failed: %v)", cause, created.ID, err)
		}
		return cause
	}

	listed, err := doPagedJSON[momentItem](ctx, publicAPI, http.MethodGet, "/moments?page=1&size=20", nil)
	if err != nil {
		return "", cleanupErr(fmt.Errorf("list public moments: %w", err))
	}
	moment, ok := findMoment(listed.Items, created.ID)
	if !ok {
		return "", cleanupErr(fmt.Errorf("public moments missing created moment %s", created.ID))
	}
	if moment.Content != content {
		return "", cleanupErr(fmt.Errorf("public moment %s returned unexpected content %q", created.ID, moment.Content))
	}

	liked, err := doJSON[map[string]bool](ctx, publicAPI, http.MethodPost, "/moments/"+created.ID+"/like", map[string]string{})
	if err != nil {
		return "", cleanupErr(fmt.Errorf("like moment %s: %w", created.ID, err))
	}
	if !liked["liked"] {
		return "", cleanupErr(fmt.Errorf("like moment %s did not return liked=true", created.ID))
	}

	reposted, err := doJSON[map[string]bool](ctx, publicAPI, http.MethodPost, "/moments/"+created.ID+"/repost", map[string]string{})
	if err != nil {
		return "", cleanupErr(fmt.Errorf("repost moment %s: %w", created.ID, err))
	}
	if !reposted["reposted"] {
		return "", cleanupErr(fmt.Errorf("repost moment %s did not return reposted=true", created.ID))
	}

	commentContent := "Smoke moment comment " + runID
	createdComment, err := doJSON[momentComment](ctx, publicAPI, http.MethodPost, "/moments/"+created.ID+"/comments", map[string]string{
		"author_name": "Smoke Visitor " + runID,
		"content":     commentContent,
	})
	if err != nil {
		return "", cleanupErr(fmt.Errorf("create moment comment for %s: %w", created.ID, err))
	}
	if strings.TrimSpace(createdComment.ID) == "" {
		return "", cleanupErr(fmt.Errorf("moment comment for %s returned empty id", created.ID))
	}

	commentLiked, err := doJSON[map[string]bool](ctx, publicAPI, http.MethodPost, "/moments/comments/"+createdComment.ID+"/like", map[string]string{})
	if err != nil {
		return "", cleanupErr(fmt.Errorf("like moment comment %s: %w", createdComment.ID, err))
	}
	if !commentLiked["liked"] {
		return "", cleanupErr(fmt.Errorf("like moment comment %s did not return liked=true", createdComment.ID))
	}

	updated, err := doPagedJSON[momentItem](ctx, publicAPI, http.MethodGet, "/moments?page=1&size=20", nil)
	if err != nil {
		return "", cleanupErr(fmt.Errorf("reload public moments: %w", err))
	}
	moment, ok = findMoment(updated.Items, created.ID)
	if !ok {
		return "", cleanupErr(fmt.Errorf("updated public moments missing moment %s", created.ID))
	}
	if !moment.Liked || moment.LikeCount < 1 {
		return "", cleanupErr(fmt.Errorf("moment %s like state not reflected in public list", created.ID))
	}
	if !moment.Reposted || moment.RepostCount < 1 {
		return "", cleanupErr(fmt.Errorf("moment %s repost state not reflected in public list", created.ID))
	}
	if moment.CommentCount < 1 {
		return "", cleanupErr(fmt.Errorf("moment %s comment count was not incremented", created.ID))
	}

	comment, ok := findMomentComment(moment.Comments, createdComment.ID)
	if !ok {
		return "", cleanupErr(fmt.Errorf("aggregated public moment %s missing comment %s", created.ID, createdComment.ID))
	}
	if comment.Content != commentContent {
		return "", cleanupErr(fmt.Errorf("aggregated comment %s returned unexpected content %q", createdComment.ID, comment.Content))
	}
	if !comment.Liked || comment.LikeCount < 1 {
		return "", cleanupErr(fmt.Errorf("aggregated comment %s like state not reflected", createdComment.ID))
	}

	if err := doNoData(ctx, adminAPI, http.MethodDelete, "/admin/moments/"+created.ID, nil); err != nil {
		return "", fmt.Errorf("cleanup moment %s: %w", created.ID, err)
	}

	return created.ID, nil
}

func newAPIClient(baseURL string, timeout time.Duration) (*apiClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &apiClient{
		baseURL: baseURL,
		client: &http.Client{
			Jar:     jar,
			Timeout: timeout,
		},
	}, nil
}

func doJSON[T any](ctx context.Context, api *apiClient, method, path string, payload any) (T, error) {
	var zero T

	requestBody, err := marshalBody(payload)
	if err != nil {
		return zero, err
	}

	req, err := http.NewRequestWithContext(ctx, method, api.baseURL+normalizePath(path), requestBody)
	if err != nil {
		return zero, fmt.Errorf("build request %s %s: %w", method, path, err)
	}
	req.Header.Set("Accept", "application/json")
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return zero, fmt.Errorf("request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("read response %s %s: %w", method, path, err)
	}

	var envelope apiEnvelope[T]
	if len(body) > 0 {
		if err := json.Unmarshal(body, &envelope); err != nil {
			return zero, fmt.Errorf("decode response %s %s: %w (body=%s)", method, path, err, strings.TrimSpace(string(body)))
		}
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return zero, fmt.Errorf("%s %s returned %d: %s", method, path, resp.StatusCode, nonEmpty(envelope.Message, strings.TrimSpace(string(body))))
	}
	if envelope.Code != 0 {
		return zero, fmt.Errorf("%s %s returned code=%d: %s", method, path, envelope.Code, envelope.Message)
	}

	return envelope.Data, nil
}

func doPagedJSON[T any](ctx context.Context, api *apiClient, method, path string, payload any) (*pagedData[T], error) {
	data, err := doJSON[pagedData[T]](ctx, api, method, path, payload)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func doNoData(ctx context.Context, api *apiClient, method, path string, payload any) error {
	_, err := doJSON[map[string]any](ctx, api, method, path, payload)
	return err
}

func expectStatus(ctx context.Context, api *apiClient, method, path string, wantStatus int) error {
	req, err := http.NewRequestWithContext(ctx, method, api.baseURL+normalizePath(path), nil)
	if err != nil {
		return err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != wantStatus {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected %d from %s %s, got %d: %s", wantStatus, method, path, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func marshalBody(payload any) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	return bytes.NewReader(buf), nil
}

func normalizeBaseURL(value string) string {
	base := strings.TrimRight(strings.TrimSpace(value), "/")
	if base == "" {
		return defaultSmokeBaseURL
	}
	if strings.HasSuffix(base, "/api/v1") {
		return base
	}
	return base + "/api/v1"
}

func normalizePath(value string) string {
	if strings.HasPrefix(value, "/") {
		return value
	}
	return "/" + value
}

func containsGuestbookMessage(items []adminGuestbookMessage, id, authorName string) bool {
	for _, item := range items {
		if item.ID == id && item.AuthorName == authorName {
			return true
		}
	}
	return false
}

func containsGuestbookID(items []guestbookMessageResponse, id string) bool {
	for _, item := range items {
		if item.ID == id {
			return true
		}
	}
	return false
}

func containsFriendApplication(items []adminFriendApplication, id, siteURL string) bool {
	for _, item := range items {
		if item.ID == id && item.SiteURL == siteURL {
			return true
		}
	}
	return false
}

func containsFriendLink(items []adminFriendLink, id, siteURL string) bool {
	for _, item := range items {
		if item.ID == id && item.URL == siteURL {
			return true
		}
	}
	return false
}

func findMoment(items []momentItem, id string) (momentItem, bool) {
	for _, item := range items {
		if item.ID == id {
			return item, true
		}
	}
	return momentItem{}, false
}

func findMomentComment(items []momentComment, id string) (momentComment, bool) {
	for _, item := range items {
		if item.ID == id {
			return item, true
		}
	}
	return momentComment{}, false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func nonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

var errUnexpectedStatus = errors.New("unexpected status")
