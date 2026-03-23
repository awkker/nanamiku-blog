package service

import (
	"context"
	"fmt"
	"strings"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuestbookService struct {
	q  *query.Queries
	db *pgxpool.Pool
}

func NewGuestbookService(db *pgxpool.Pool) *GuestbookService {
	return &GuestbookService{q: query.New(db), db: db}
}

type GuestbookMessageItem struct {
	ID            uuid.UUID              `json:"id"`
	ParentID      *uuid.UUID             `json:"parent_id,omitempty"`
	AuthorName    string                 `json:"author_name"`
	AuthorWebsite string                 `json:"author_website,omitempty"`
	IsAuthor      bool                   `json:"is_author"`
	AuthorAvatar  string                 `json:"author_avatar_url,omitempty"`
	Content       string                 `json:"content"`
	VoteScore     int32                  `json:"vote_score"`
	CreatedAt     string                 `json:"created_at"`
	Replies       []GuestbookMessageItem `json:"replies,omitempty"`
	MyVote        *string                `json:"my_vote,omitempty"`
}

type GuestbookAuthorIdentity struct {
	AdminID     uuid.UUID
	DisplayName string
	AvatarURL   string
}

type guestbookAuthorMark struct {
	AvatarURL string
}

const upsertGuestbookAuthorMarkSQL = `
INSERT INTO guestbook_author_marks (message_id, admin_id, admin_avatar_url)
VALUES ($1, $2, $3)
ON CONFLICT (message_id) DO UPDATE
SET admin_id = EXCLUDED.admin_id,
    admin_avatar_url = EXCLUDED.admin_avatar_url
`

const listGuestbookAuthorMarksSQL = `
SELECT message_id, admin_avatar_url
FROM guestbook_author_marks
WHERE message_id = ANY($1::uuid[])
`

func (s *GuestbookService) ListMessages(ctx context.Context, sortBy string, page, size int, visitorID uuid.UUID) ([]GuestbookMessageItem, int64, error) {
	total, err := s.q.CountGuestbookMessages(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count messages: %w", err)
	}

	rows, err := s.q.ListGuestbookMessages(ctx, query.ListGuestbookMessagesParams{
		SortBy: sortBy,
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list messages: %w", err)
	}

	msgIDs := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		msgIDs[i] = r.ID
	}

	var voteMap map[uuid.UUID]string
	if visitorID != uuid.Nil && len(msgIDs) > 0 {
		voteMap = s.getVoteMap(ctx, visitorID, msgIDs)
	}

	items := make([]GuestbookMessageItem, 0, len(rows))
	allIDs := make([]uuid.UUID, 0, len(rows)*2)
	for _, r := range rows {
		item := GuestbookMessageItem{
			ID:            r.ID,
			AuthorName:    r.AuthorName,
			AuthorWebsite: r.AuthorWebsite,
			Content:       r.Content,
			VoteScore:     r.VoteScore,
			CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if v, ok := voteMap[r.ID]; ok {
			item.MyVote = &v
		}
		allIDs = append(allIDs, r.ID)

		replies, _ := s.q.ListGuestbookReplies(ctx, pgtype.UUID{Bytes: r.ID, Valid: true})
		replyItems := make([]GuestbookMessageItem, 0, len(replies))
		for _, rp := range replies {
			ri := GuestbookMessageItem{
				ID:            rp.ID,
				AuthorName:    rp.AuthorName,
				AuthorWebsite: rp.AuthorWebsite,
				Content:       rp.Content,
				VoteScore:     rp.VoteScore,
				CreatedAt:     rp.CreatedAt.Format("2006-01-02T15:04:05Z"),
			}
			allIDs = append(allIDs, rp.ID)
			replyItems = append(replyItems, ri)
		}
		item.Replies = replyItems
		items = append(items, item)
	}

	authorMarkMap := s.getAuthorMarkMap(ctx, allIDs)
	applyGuestbookAuthorMarks(items, authorMarkMap)

	return items, total, nil
}

func (s *GuestbookService) CreateMessage(
	ctx context.Context,
	parentID *uuid.UUID,
	authorName, authorWebsite, content, ipHash, uaHash string,
	authorIdentity *GuestbookAuthorIdentity,
) (*GuestbookMessageItem, error) {
	var pid pgtype.UUID
	if parentID != nil {
		pid = pgtype.UUID{Bytes: *parentID, Valid: true}
	}

	effectiveAuthorName := strings.TrimSpace(authorName)
	effectiveAuthorAvatar := ""
	if authorIdentity != nil {
		displayName := strings.TrimSpace(authorIdentity.DisplayName)
		if displayName != "" {
			effectiveAuthorName = displayName
		}
		effectiveAuthorAvatar = strings.TrimSpace(authorIdentity.AvatarURL)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create guestbook message tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	qtx := s.q.WithTx(tx)
	row, err := qtx.CreateGuestbookMessage(ctx, query.CreateGuestbookMessageParams{
		ParentID:      pid,
		AuthorName:    effectiveAuthorName,
		AuthorWebsite: authorWebsite,
		Content:       content,
		IpHash:        ipHash,
		UaHash:        uaHash,
	})
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	if authorIdentity != nil && authorIdentity.AdminID != uuid.Nil {
		if err := qtx.ApproveGuestbookMessage(ctx, query.ApproveGuestbookMessageParams{
			ID:         row.ID,
			ReviewedBy: pgtype.UUID{Bytes: authorIdentity.AdminID, Valid: true},
		}); err != nil {
			return nil, fmt.Errorf("auto approve author message: %w", err)
		}

		if _, err := tx.Exec(ctx, upsertGuestbookAuthorMarkSQL, row.ID, authorIdentity.AdminID, effectiveAuthorAvatar); err != nil {
			return nil, fmt.Errorf("mark author guestbook message: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create guestbook message tx: %w", err)
	}
	committed = true

	isAuthor := authorIdentity != nil && authorIdentity.AdminID != uuid.Nil
	return &GuestbookMessageItem{
		ID:            row.ID,
		AuthorName:    effectiveAuthorName,
		AuthorWebsite: authorWebsite,
		IsAuthor:      isAuthor,
		AuthorAvatar:  effectiveAuthorAvatar,
		Content:       content,
		VoteScore:     0,
		CreatedAt:     row.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GuestbookService) Vote(ctx context.Context, messageID, visitorID uuid.UUID, vote string) error {
	if vote == "" {
		_ = s.q.DeleteGuestbookVote(ctx, query.DeleteGuestbookVoteParams{
			MessageID: messageID,
			VisitorID: visitorID,
		})
	} else {
		vt := query.VoteType(vote)
		_ = s.q.UpsertGuestbookVote(ctx, query.UpsertGuestbookVoteParams{
			MessageID: messageID,
			VisitorID: visitorID,
			Vote:      vt,
		})
	}

	return s.q.RecalcGuestbookVoteScore(ctx, messageID)
}

func (s *GuestbookService) getVoteMap(ctx context.Context, visitorID uuid.UUID, msgIDs []uuid.UUID) map[uuid.UUID]string {
	m := make(map[uuid.UUID]string)
	votes, err := s.q.GetVisitorVotesForMessages(ctx, query.GetVisitorVotesForMessagesParams{
		VisitorID: visitorID,
		Column2:   msgIDs,
	})
	if err != nil {
		return m
	}
	for _, v := range votes {
		m[v.MessageID] = string(v.Vote)
	}
	return m
}

func (s *GuestbookService) getAuthorMarkMap(ctx context.Context, msgIDs []uuid.UUID) map[uuid.UUID]guestbookAuthorMark {
	m := make(map[uuid.UUID]guestbookAuthorMark)
	if len(msgIDs) == 0 {
		return m
	}

	rows, err := s.db.Query(ctx, listGuestbookAuthorMarksSQL, msgIDs)
	if err != nil {
		return m
	}
	defer rows.Close()

	for rows.Next() {
		var messageID uuid.UUID
		var avatarURL string
		if err := rows.Scan(&messageID, &avatarURL); err != nil {
			return m
		}
		m[messageID] = guestbookAuthorMark{AvatarURL: avatarURL}
	}
	return m
}

func applyGuestbookAuthorMarks(items []GuestbookMessageItem, markMap map[uuid.UUID]guestbookAuthorMark) {
	if len(markMap) == 0 || len(items) == 0 {
		return
	}
	for idx := range items {
		applyGuestbookAuthorMark(&items[idx], markMap)
	}
}

func applyGuestbookAuthorMark(item *GuestbookMessageItem, markMap map[uuid.UUID]guestbookAuthorMark) {
	if mark, ok := markMap[item.ID]; ok {
		item.IsAuthor = true
		item.AuthorAvatar = mark.AvatarURL
	}
	for idx := range item.Replies {
		applyGuestbookAuthorMark(&item.Replies[idx], markMap)
	}
}
