package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ModerationService struct {
	q   *query.Queries
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewModerationService(db *pgxpool.Pool, rdb *redis.Client) *ModerationService {
	return &ModerationService{q: query.New(db), db: db, rdb: rdb}
}

func (s *ModerationService) ListSensitiveWords(ctx context.Context) ([]string, error) {
	return s.q.ListSensitiveWords(ctx)
}

func (s *ModerationService) CheckBlocked(ctx context.Context, ipHash string) (bool, error) {
	cnt, err := s.q.CheckBlockedIP(ctx, ipHash)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (s *ModerationService) FindSensitiveWord(ctx context.Context, texts ...string) (string, error) {
	words, err := s.q.ListSensitiveWords(ctx)
	if err != nil {
		return "", fmt.Errorf("list sensitive words: %w", err)
	}
	if len(words) == 0 || len(texts) == 0 {
		return "", nil
	}

	normalizedTexts := make([]string, 0, len(texts))
	for _, text := range texts {
		t := strings.ToLower(strings.TrimSpace(text))
		if t != "" {
			normalizedTexts = append(normalizedTexts, t)
		}
	}
	if len(normalizedTexts) == 0 {
		return "", nil
	}

	for _, word := range words {
		token := strings.ToLower(strings.TrimSpace(word))
		if token == "" {
			continue
		}
		for _, text := range normalizedTexts {
			if strings.Contains(text, token) {
				return word, nil
			}
		}
	}
	return "", nil
}

var defaultRateLimitRules = []string{
	"analytics:collect",
	"friend:apply",
	"gb:create",
	"gb:vote",
	"login",
	"mt:clike",
	"mt:comment",
	"mt:create",
	"mt:like",
	"mt:repost",
	"post:comment",
	"post:like",
}

type RateLimitRuleMetric struct {
	Rule    string `json:"rule"`
	Allowed int64  `json:"allowed"`
	Blocked int64  `json:"blocked"`
	Total   int64  `json:"total"`
}

type RateLimitTrendPoint struct {
	Bucket  string `json:"bucket"`
	Allowed int64  `json:"allowed"`
	Blocked int64  `json:"blocked"`
}

type RateLimitMetrics struct {
	WindowMinutes int                   `json:"window_minutes"`
	TotalAllowed  int64                 `json:"total_allowed"`
	TotalBlocked  int64                 `json:"total_blocked"`
	Rules         []RateLimitRuleMetric `json:"rules"`
	Trend         []RateLimitTrendPoint `json:"trend"`
}

func (s *ModerationService) GetRateLimitMetrics(ctx context.Context, minutes int) (*RateLimitMetrics, error) {
	if s.rdb == nil {
		return nil, fmt.Errorf("redis not configured")
	}
	if minutes <= 0 {
		minutes = 60
	}
	if minutes > 24*60 {
		minutes = 24 * 60
	}

	rules, err := s.rdb.SMembers(ctx, "rlm:rules").Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("read rate-limit rules: %w", err)
	}
	rules = normalizeRateLimitRules(rules)
	if len(rules) == 0 {
		rules = append([]string{}, defaultRateLimitRules...)
	}

	now := time.Now().UTC().Truncate(time.Minute)
	start := now.Add(-time.Duration(minutes-1) * time.Minute)
	buckets := make([]time.Time, 0, minutes)
	for i := 0; i < minutes; i++ {
		buckets = append(buckets, start.Add(time.Duration(i)*time.Minute))
	}

	keys := make([]string, 0, len(rules)*len(buckets)*2)
	for _, rule := range rules {
		for _, bucket := range buckets {
			b := bucket.Format("200601021504")
			keys = append(keys,
				fmt.Sprintf("rlm:%s:%s:allow", rule, b),
				fmt.Sprintf("rlm:%s:%s:block", rule, b),
			)
		}
	}

	values := make([]interface{}, 0)
	if len(keys) > 0 {
		values, err = s.rdb.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, fmt.Errorf("read rate-limit metrics: %w", err)
		}
	}

	trend := make([]RateLimitTrendPoint, len(buckets))
	for i, bucket := range buckets {
		trend[i] = RateLimitTrendPoint{Bucket: bucket.Format("2006-01-02 15:04")}
	}

	ruleMetrics := make([]RateLimitRuleMetric, 0, len(rules))
	var totalAllowed int64
	var totalBlocked int64

	idx := 0
	for _, rule := range rules {
		item := RateLimitRuleMetric{Rule: rule}
		for i := range buckets {
			allowed := redisValueToInt64(values[idx])
			idx++
			blocked := redisValueToInt64(values[idx])
			idx++

			item.Allowed += allowed
			item.Blocked += blocked
			trend[i].Allowed += allowed
			trend[i].Blocked += blocked
		}
		item.Total = item.Allowed + item.Blocked
		totalAllowed += item.Allowed
		totalBlocked += item.Blocked
		ruleMetrics = append(ruleMetrics, item)
	}

	return &RateLimitMetrics{
		WindowMinutes: minutes,
		TotalAllowed:  totalAllowed,
		TotalBlocked:  totalBlocked,
		Rules:         ruleMetrics,
		Trend:         trend,
	}, nil
}

func normalizeRateLimitRules(rules []string) []string {
	seen := make(map[string]struct{}, len(rules))
	out := make([]string, 0, len(rules))
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		if _, ok := seen[rule]; ok {
			continue
		}
		seen[rule] = struct{}{}
		out = append(out, rule)
	}
	sort.Strings(out)
	return out
}

func redisValueToInt64(v interface{}) int64 {
	switch t := v.(type) {
	case nil:
		return 0
	case int64:
		return t
	case string:
		n, _ := strconv.ParseInt(t, 10, 64)
		return n
	case []byte:
		n, _ := strconv.ParseInt(string(t), 10, 64)
		return n
	default:
		return 0
	}
}

func (s *ModerationService) LogAudit(ctx context.Context, adminID uuid.UUID, action, targetType, targetID string, detail interface{}, ip string) error {
	detailJSON, _ := json.Marshal(detail)
	return s.q.CreateAuditLog(ctx, query.CreateAuditLogParams{
		AdminID:    pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detailJSON,
		Ip:         ip,
	})
}

type AuditLogItem struct {
	ID            uuid.UUID       `json:"id"`
	Action        string          `json:"action"`
	TargetType    string          `json:"target_type"`
	TargetID      string          `json:"target_id"`
	Detail        json.RawMessage `json:"detail"`
	IP            string          `json:"ip"`
	AdminUsername string          `json:"admin_username,omitempty"`
	CreatedAt     string          `json:"created_at"`
}

func (s *ModerationService) ListAuditLogs(ctx context.Context, page, size int) ([]AuditLogItem, error) {
	rows, err := s.q.ListAuditLogs(ctx, query.ListAuditLogsParams{
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}

	items := make([]AuditLogItem, 0, len(rows))
	for _, r := range rows {
		var username string
		if r.AdminUsername.Valid {
			username = r.AdminUsername.String
		}
		items = append(items, AuditLogItem{
			ID:            r.ID,
			Action:        r.Action,
			TargetType:    r.TargetType,
			TargetID:      r.TargetID,
			Detail:        r.Detail,
			IP:            r.Ip,
			AdminUsername: username,
			CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return items, nil
}

type AdminCommentItem struct {
	ID          uuid.UUID `json:"id"`
	PostID      uuid.UUID `json:"post_id"`
	PostTitle   string    `json:"post_title,omitempty"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	IPHash      string    `json:"ip_hash"`
	CreatedAt   string    `json:"created_at"`
}

func (s *ModerationService) ListComments(ctx context.Context, status string, page, size int) ([]AdminCommentItem, int64, error) {
	var statusParam query.NullModerationStatus
	if status != "" {
		statusParam = query.NullModerationStatus{
			ModerationStatus: query.ModerationStatus(status),
			Valid:            true,
		}
	}

	total, err := s.q.CountAdminComments(ctx, statusParam)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListAdminComments(ctx, query.ListAdminCommentsParams{
		Status: statusParam,
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]AdminCommentItem, 0, len(rows))
	for _, r := range rows {
		var postTitle string
		if r.PostTitle.Valid {
			postTitle = r.PostTitle.String
		}
		items = append(items, AdminCommentItem{
			ID:          r.ID,
			PostID:      r.PostID,
			PostTitle:   postTitle,
			AuthorName:  r.AuthorName,
			AuthorEmail: r.AuthorEmail,
			Content:     r.Content,
			Status:      string(r.Status),
			IPHash:      r.IpHash,
			CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return items, total, nil
}

func (s *ModerationService) ApproveComment(ctx context.Context, commentID, adminID uuid.UUID) error {
	return s.q.ApproveComment(ctx, query.ApproveCommentParams{ID: commentID, ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: true}})
}

func (s *ModerationService) RejectComment(ctx context.Context, commentID, adminID uuid.UUID) error {
	return s.q.RejectComment(ctx, query.RejectCommentParams{ID: commentID, ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: true}})
}

func (s *ModerationService) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	return s.q.DeleteComment(ctx, commentID)
}

type AdminGuestbookMessageItem struct {
	ID               uuid.UUID  `json:"id"`
	ParentID         *uuid.UUID `json:"parent_id,omitempty"`
	ParentAuthorName string     `json:"parent_author_name,omitempty"`
	AuthorName       string     `json:"author_name"`
	AuthorWebsite    string     `json:"author_website,omitempty"`
	Content          string     `json:"content"`
	Status           string     `json:"status"`
	IPHash           string     `json:"ip_hash"`
	VoteScore        int32      `json:"vote_score"`
	CreatedAt        string     `json:"created_at"`
}

func (s *ModerationService) ListGuestbookMessages(ctx context.Context, status string, page, size int) ([]AdminGuestbookMessageItem, int64, error) {
	var statusParam query.NullModerationStatus
	if status != "" {
		statusParam = query.NullModerationStatus{
			ModerationStatus: query.ModerationStatus(status),
			Valid:            true,
		}
	}

	total, err := s.q.CountAdminGuestbookMessages(ctx, statusParam)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListAdminGuestbookMessages(ctx, query.ListAdminGuestbookMessagesParams{
		Status: statusParam,
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]AdminGuestbookMessageItem, 0, len(rows))
	for _, r := range rows {
		item := AdminGuestbookMessageItem{
			ID:            r.ID,
			AuthorName:    r.AuthorName,
			AuthorWebsite: r.AuthorWebsite,
			Content:       r.Content,
			Status:        string(r.Status),
			IPHash:        r.IpHash,
			VoteScore:     r.VoteScore,
			CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if r.ParentID.Valid {
			pid := uuid.UUID(r.ParentID.Bytes)
			item.ParentID = &pid
		}
		if r.ParentAuthorName.Valid {
			item.ParentAuthorName = r.ParentAuthorName.String
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (s *ModerationService) ApproveGuestbookMessage(ctx context.Context, messageID, adminID uuid.UUID) error {
	return s.q.ApproveGuestbookMessage(ctx, query.ApproveGuestbookMessageParams{
		ID:         messageID,
		ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: true},
	})
}

func (s *ModerationService) RejectGuestbookMessage(ctx context.Context, messageID, adminID uuid.UUID) error {
	return s.q.RejectGuestbookMessage(ctx, query.RejectGuestbookMessageParams{
		ID:         messageID,
		ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: true},
	})
}

func (s *ModerationService) DeleteGuestbookMessage(ctx context.Context, messageID uuid.UUID) error {
	return s.q.DeleteGuestbookMessage(ctx, messageID)
}

type AdminFriendItem struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	URL          string    `json:"url"`
	Domain       string    `json:"domain"`
	AvatarURL    string    `json:"avatar_url"`
	Status       string    `json:"status"`
	HealthStatus string    `json:"health_status"`
	SortOrder    int32     `json:"sort_order"`
	CreatedAt    string    `json:"created_at"`
}

func (s *ModerationService) ListAdminFriends(ctx context.Context, page, size int) ([]AdminFriendItem, int64, error) {
	total, err := s.q.CountAdminFriendLinks(ctx)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListAdminFriendLinks(ctx, query.ListAdminFriendLinksParams{
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]AdminFriendItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, AdminFriendItem{
			ID:           r.ID,
			Name:         r.Name,
			Description:  r.Description,
			URL:          r.Url,
			Domain:       r.Domain,
			AvatarURL:    r.AvatarUrl,
			Status:       string(r.Status),
			HealthStatus: string(r.HealthStatus),
			SortOrder:    r.SortOrder,
			CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return items, total, nil
}

func (s *ModerationService) CreateFriend(ctx context.Context, name, description, url, domain, avatarURL string, sortOrder int32, adminID uuid.UUID) (uuid.UUID, error) {
	normalized, err := normalizeFriendLinkInput(FriendLinkInput{
		Name:        name,
		Description: description,
		URL:         url,
		Domain:      domain,
		AvatarURL:   avatarURL,
	})
	if err != nil {
		return uuid.Nil, err
	}

	row, err := s.q.CreateFriendLink(ctx, query.CreateFriendLinkParams{
		Name:        normalized.Name,
		Description: normalized.Description,
		Url:         normalized.URL,
		Domain:      normalized.Domain,
		AvatarUrl:   normalized.AvatarURL,
		SortOrder:   sortOrder,
		ReviewedBy:  pgtype.UUID{Bytes: adminID, Valid: true},
	})
	if err != nil {
		return uuid.Nil, err
	}
	return row.ID, nil
}

func (s *ModerationService) UpdateFriend(ctx context.Context, id uuid.UUID, name, description, url, domain, avatarURL string, sortOrder int32) error {
	normalized, err := normalizeFriendLinkInput(FriendLinkInput{
		Name:        name,
		Description: description,
		URL:         url,
		Domain:      domain,
		AvatarURL:   avatarURL,
	})
	if err != nil {
		return err
	}

	return s.q.UpdateFriendLink(ctx, query.UpdateFriendLinkParams{
		ID:          id,
		Name:        normalized.Name,
		Description: normalized.Description,
		Url:         normalized.URL,
		Domain:      normalized.Domain,
		AvatarUrl:   normalized.AvatarURL,
		SortOrder:   sortOrder,
	})
}

func (s *ModerationService) DeleteFriend(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteFriendLink(ctx, id)
}

func (s *ModerationService) ApproveFriend(ctx context.Context, id uuid.UUID, adminID uuid.UUID) error {
	return s.q.ApproveFriendLink(ctx, query.ApproveFriendLinkParams{
		ID:         id,
		ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
}

func (s *ModerationService) RejectFriend(ctx context.Context, id uuid.UUID, adminID uuid.UUID) error {
	return s.q.RejectFriendLink(ctx, query.RejectFriendLinkParams{
		ID:         id,
		ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
}

type AdminFriendApplicationItem struct {
	ID           uuid.UUID `json:"id"`
	SiteName     string    `json:"site_name"`
	SiteURL      string    `json:"site_url"`
	AvatarURL    string    `json:"avatar_url"`
	Description  string    `json:"description"`
	ContactEmail string    `json:"contact_email"`
	ContactNote  string    `json:"contact_note"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"created_at"`
	ReviewedAt   string    `json:"reviewed_at,omitempty"`
	ReviewNote   string    `json:"review_note"`
}

func (s *ModerationService) ListFriendApplications(ctx context.Context, status string, page, size int) ([]AdminFriendApplicationItem, int64, error) {
	var statusParam query.NullFriendLinkStatus
	if status == string(query.FriendLinkStatusPending) || status == string(query.FriendLinkStatusApproved) || status == string(query.FriendLinkStatusRejected) {
		statusParam = query.NullFriendLinkStatus{
			FriendLinkStatus: query.FriendLinkStatus(status),
			Valid:            true,
		}
	}

	total, err := s.q.CountAdminFriendLinkApplications(ctx, statusParam)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListAdminFriendLinkApplications(ctx, query.ListAdminFriendLinkApplicationsParams{
		Status: statusParam,
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]AdminFriendApplicationItem, 0, len(rows))
	for _, r := range rows {
		item := AdminFriendApplicationItem{
			ID:           r.ID,
			SiteName:     r.SiteName,
			SiteURL:      r.SiteUrl,
			AvatarURL:    r.AvatarUrl,
			Description:  r.Description,
			ContactEmail: r.ContactEmail,
			ContactNote:  r.ContactNote,
			Status:       string(r.Status),
			CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
			ReviewNote:   r.ReviewNote,
		}
		if r.ReviewedAt.Valid {
			item.ReviewedAt = r.ReviewedAt.Time.Format("2006-01-02T15:04:05Z")
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (s *ModerationService) ApproveFriendApplication(ctx context.Context, applicationID, adminID uuid.UUID, sortOrder int32, reviewNote string) (uuid.UUID, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin friend application tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	application, err := qtx.GetFriendLinkApplicationByID(ctx, applicationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrFriendApplicationNotFound
		}
		return uuid.Nil, fmt.Errorf("get friend application: %w", err)
	}
	if application.Status != query.FriendLinkStatusPending {
		return uuid.Nil, ErrFriendApplicationProcessed
	}

	normalized, err := normalizeFriendLinkInput(FriendLinkInput{
		Name:        application.SiteName,
		Description: application.Description,
		URL:         application.SiteUrl,
		AvatarURL:   application.AvatarUrl,
	})
	if err != nil {
		return uuid.Nil, err
	}

	exists, err := qtx.FriendLinkExistsByURL(ctx, normalized.URL)
	if err != nil {
		return uuid.Nil, fmt.Errorf("check existing friend link: %w", err)
	}
	if exists {
		return uuid.Nil, ErrFriendApplicationDuplicate
	}

	friendRow, err := qtx.CreateFriendLink(ctx, query.CreateFriendLinkParams{
		Name:        normalized.Name,
		Description: normalized.Description,
		Url:         normalized.URL,
		Domain:      normalized.Domain,
		AvatarUrl:   normalized.AvatarURL,
		SortOrder:   sortOrder,
		ReviewedBy:  pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create friend link from application: %w", err)
	}

	if err := qtx.ApproveFriendLinkApplication(ctx, query.ApproveFriendLinkApplicationParams{
		ID:         applicationID,
		ReviewNote: strings.TrimSpace(reviewNote),
		ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	}); err != nil {
		return uuid.Nil, fmt.Errorf("approve friend application: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit friend application approval: %w", err)
	}

	return friendRow.ID, nil
}

func (s *ModerationService) RejectFriendApplication(ctx context.Context, applicationID, adminID uuid.UUID, reviewNote string) error {
	application, err := s.q.GetFriendLinkApplicationByID(ctx, applicationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrFriendApplicationNotFound
		}
		return fmt.Errorf("get friend application: %w", err)
	}
	if application.Status != query.FriendLinkStatusPending {
		return ErrFriendApplicationProcessed
	}

	return s.q.RejectFriendLinkApplication(ctx, query.RejectFriendLinkApplicationParams{
		ID:         applicationID,
		ReviewNote: strings.TrimSpace(reviewNote),
		ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
}
