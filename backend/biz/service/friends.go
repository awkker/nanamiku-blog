package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FriendsService struct {
	q  *query.Queries
	db *pgxpool.Pool
}

func NewFriendsService(db *pgxpool.Pool) *FriendsService {
	return &FriendsService{q: query.New(db), db: db}
}

type FriendLinkItem struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	URL          string    `json:"url"`
	AvatarURL    string    `json:"avatar_url"`
	HealthStatus string    `json:"health_status"`
}

func (s *FriendsService) ListApproved(ctx context.Context) ([]FriendLinkItem, error) {
	rows, err := s.q.ListApprovedFriendLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list friend links: %w", err)
	}

	items := make([]FriendLinkItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, FriendLinkItem{
			ID:           r.ID,
			Name:         r.Name,
			Description:  r.Description,
			URL:          r.Url,
			AvatarURL:    r.AvatarUrl,
			HealthStatus: string(r.HealthStatus),
		})
	}
	return items, nil
}

type FriendApplicationSubmission struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
}

func (s *FriendsService) SubmitApplication(ctx context.Context, input FriendApplicationInput) (*FriendApplicationSubmission, error) {
	normalized, err := normalizeFriendApplicationInput(input)
	if err != nil {
		return nil, err
	}

	exists, err := s.q.FriendLinkExistsByURL(ctx, normalized.SiteURL)
	if err != nil {
		return nil, fmt.Errorf("check existing friend link: %w", err)
	}
	if exists {
		return nil, ErrFriendApplicationDuplicate
	}

	pendingExists, err := s.q.PendingFriendLinkApplicationExistsByURL(ctx, normalized.SiteURL)
	if err != nil {
		return nil, fmt.Errorf("check pending friend application: %w", err)
	}
	if pendingExists {
		return nil, ErrFriendApplicationDuplicate
	}

	row, err := s.q.CreateFriendLinkApplication(ctx, query.CreateFriendLinkApplicationParams{
		SiteName:     normalized.SiteName,
		SiteUrl:      normalized.SiteURL,
		AvatarUrl:    normalized.AvatarURL,
		Description:  normalized.Description,
		ContactEmail: normalized.ContactEmail,
		ContactNote:  normalized.ContactNote,
	})
	if err != nil {
		return nil, fmt.Errorf("create friend application: %w", err)
	}

	return &FriendApplicationSubmission{
		ID:        row.ID,
		Status:    string(row.Status),
		CreatedAt: row.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *FriendsService) RunHealthChecks(ctx context.Context) {
	links, err := s.q.ListFriendLinksForHealthCheck(ctx)
	if err != nil {
		slog.Error("health check: list links failed", "error", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, link := range links {
		go func(id uuid.UUID, url string) {
			start := time.Now()
			status := query.HealthStatusOk
			httpStatus := 0

			resp, err := client.Get(url)
			latency := time.Since(start).Milliseconds()

			if err != nil {
				status = query.HealthStatusDown
			} else {
				httpStatus = resp.StatusCode
				resp.Body.Close()
				if resp.StatusCode >= 400 {
					status = query.HealthStatusDown
				}
			}

			bgCtx := context.Background()
			_ = s.q.CreateFriendLinkHealthLog(bgCtx, query.CreateFriendLinkHealthLogParams{
				FriendLinkID: id,
				HttpStatus:   int32(httpStatus),
				LatencyMs:    int32(latency),
				Result:       status,
			})
			_ = s.q.UpdateFriendLinkHealth(bgCtx, query.UpdateFriendLinkHealthParams{
				ID:           id,
				HealthStatus: status,
			})

			slog.Info("health check done", "link_id", id, "url", url, "status", status, "latency_ms", latency)
		}(link.ID, link.Url)
	}
}
