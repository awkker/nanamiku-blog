package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nanamiku-blog/backend/query"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	dashStatsCacheKey = "dash:stats"
	dashStatsTTL      = 2 * time.Minute
	dashTrendTTL      = 5 * time.Minute
)

type DashboardService struct {
	q           *query.Queries
	db          *pgxpool.Pool
	rdb         *redis.Client
	geoResolver *GeoIPResolver
}

func NewDashboardService(db *pgxpool.Pool, rdb *redis.Client, geoResolver *GeoIPResolver) *DashboardService {
	return &DashboardService{q: query.New(db), db: db, rdb: rdb, geoResolver: geoResolver}
}

type DashboardStats struct {
	TotalPosts   int64 `json:"total_posts"`
	TotalLikes   int64 `json:"total_likes"`
	PendingCount int64 `json:"pending_comments"`
	FriendCount  int64 `json:"friend_count"`
	DraftCount   int64 `json:"draft_count"`
}

type TrendPoint struct {
	Day   string `json:"day"`
	Value int64  `json:"value"`
}

type ViewTrendPoint struct {
	Day string `json:"day"`
	PV  int64  `json:"pv"`
	UV  int64  `json:"uv"`
}

func (s *DashboardService) GetStats(ctx context.Context) (*DashboardStats, error) {
	if cached, ok := s.readCache(ctx, dashStatsCacheKey, new(DashboardStats)); ok {
		return cached.(*DashboardStats), nil
	}

	totalPosts, err := s.q.GetTotalPostCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("total posts: %w", err)
	}
	totalLikes, err := s.q.GetTotalLikeCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("total likes: %w", err)
	}
	pendingComments, err := s.q.CountPendingComments(ctx)
	if err != nil {
		return nil, fmt.Errorf("pending comments: %w", err)
	}
	pendingGuestbook, err := s.q.CountAdminGuestbookMessages(ctx, query.NullModerationStatus{
		ModerationStatus: query.ModerationStatusPending,
		Valid:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("pending guestbook messages: %w", err)
	}
	friendCount, err := s.q.CountApprovedFriendLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("friend count: %w", err)
	}
	draftCount, err := s.q.CountDraftPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("draft count: %w", err)
	}

	stats := &DashboardStats{
		TotalPosts:   totalPosts,
		TotalLikes:   totalLikes,
		PendingCount: pendingComments + pendingGuestbook,
		FriendCount:  friendCount,
		DraftCount:   draftCount,
	}
	s.writeCache(ctx, dashStatsCacheKey, stats, dashStatsTTL)
	return stats, nil
}

func (s *DashboardService) GetViewTrend(ctx context.Context, days int) ([]ViewTrendPoint, error) {
	cacheKey := fmt.Sprintf("dash:view_trend:%d", days)
	var cached []ViewTrendPoint
	if v, ok := s.readCache(ctx, cacheKey, &cached); ok {
		return *v.(*[]ViewTrendPoint), nil
	}

	since := time.Now().AddDate(0, 0, -days)
	rows, err := s.q.GetDailyViewTrend(ctx, pgtype.Date{Time: since, Valid: true})
	if err != nil {
		return nil, err
	}
	points := make([]ViewTrendPoint, 0, len(rows))
	for _, r := range rows {
		points = append(points, ViewTrendPoint{
			Day: r.Day.Time.Format("2006-01-02"),
			PV:  r.Pv,
			UV:  r.Uv,
		})
	}
	s.writeCache(ctx, cacheKey, points, dashTrendTTL)
	return points, nil
}

func (s *DashboardService) GetCommentTrend(ctx context.Context, days int) ([]TrendPoint, error) {
	cacheKey := fmt.Sprintf("dash:comment_trend:%d", days)
	var cached []TrendPoint
	if v, ok := s.readCache(ctx, cacheKey, &cached); ok {
		return *v.(*[]TrendPoint), nil
	}

	since := time.Now().AddDate(0, 0, -days)
	rows, err := s.q.GetDailyCommentTrend(ctx, since)
	if err != nil {
		return nil, err
	}
	points := make([]TrendPoint, 0, len(rows))
	for _, r := range rows {
		points = append(points, TrendPoint{
			Day:   r.Day.Time.Format("2006-01-02"),
			Value: r.Total,
		})
	}
	s.writeCache(ctx, cacheKey, points, dashTrendTTL)
	return points, nil
}

func (s *DashboardService) GetLikeTrend(ctx context.Context, days int) ([]TrendPoint, error) {
	cacheKey := fmt.Sprintf("dash:like_trend:%d", days)
	var cached []TrendPoint
	if v, ok := s.readCache(ctx, cacheKey, &cached); ok {
		return *v.(*[]TrendPoint), nil
	}

	since := time.Now().AddDate(0, 0, -days)
	rows, err := s.q.GetDailyLikeTrend(ctx, since)
	if err != nil {
		return nil, err
	}
	points := make([]TrendPoint, 0, len(rows))
	for _, r := range rows {
		points = append(points, TrendPoint{
			Day:   r.Day.Time.Format("2006-01-02"),
			Value: r.Total,
		})
	}
	s.writeCache(ctx, cacheKey, points, dashTrendTTL)
	return points, nil
}

func (s *DashboardService) readCache(ctx context.Context, key string, target interface{}) (interface{}, bool) {
	if s.rdb == nil {
		return nil, false
	}
	data, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	if err := json.Unmarshal(data, target); err != nil {
		return nil, false
	}
	return target, true
}

func (s *DashboardService) writeCache(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	if s.rdb == nil {
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = s.rdb.Set(ctx, key, data, ttl).Err()
}
