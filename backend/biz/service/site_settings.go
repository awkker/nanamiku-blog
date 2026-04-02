package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	siteSettingsFooterKey = "footer"
	maxFooterCustomTexts  = 8
)

type SiteSettingsService struct {
	q  *query.Queries
	db *pgxpool.Pool
}

func NewSiteSettingsService(db *pgxpool.Pool) *SiteSettingsService {
	return &SiteSettingsService{q: query.New(db), db: db}
}

type FooterSettings struct {
	ICPText     string   `json:"icp_text"`
	ICPLink     string   `json:"icp_link"`
	CustomTexts []string `json:"custom_texts"`
}

type footerSettingsDoc struct {
	ICPText     string   `json:"icp_text"`
	ICPLink     string   `json:"icp_link"`
	CustomTexts []string `json:"custom_texts"`
}

func (s *SiteSettingsService) GetFooterSettings(ctx context.Context) (*FooterSettings, error) {
	row, err := s.q.GetSiteSetting(ctx, siteSettingsFooterKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get footer settings: %w", err)
	}

	settings, err := decodeFooterSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode footer settings: %w", err)
	}

	return settings, nil
}

func (s *SiteSettingsService) SaveFooterSettings(ctx context.Context, input FooterSettings, adminID uuid.UUID) (*FooterSettings, error) {
	normalized := normalizeFooterSettings(input)
	payload, err := json.Marshal(footerSettingsDoc{
		ICPText:     normalized.ICPText,
		ICPLink:     normalized.ICPLink,
		CustomTexts: normalized.CustomTexts,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal footer settings: %w", err)
	}

	row, err := s.q.UpsertSiteSetting(ctx, query.UpsertSiteSettingParams{
		Key:       siteSettingsFooterKey,
		Value:     payload,
		UpdatedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
	if err != nil {
		return nil, fmt.Errorf("save footer settings: %w", err)
	}

	settings, err := decodeFooterSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode saved footer settings: %w", err)
	}

	return settings, nil
}

func decodeFooterSettings(raw json.RawMessage) (*FooterSettings, error) {
	if len(raw) == 0 {
		return &FooterSettings{CustomTexts: []string{}}, nil
	}

	var doc footerSettingsDoc
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	normalized := normalizeFooterSettings(FooterSettings{
		ICPText:     doc.ICPText,
		ICPLink:     doc.ICPLink,
		CustomTexts: doc.CustomTexts,
	})

	return &normalized, nil
}

func normalizeFooterSettings(input FooterSettings) FooterSettings {
	return FooterSettings{
		ICPText:     strings.TrimSpace(input.ICPText),
		ICPLink:     strings.TrimSpace(input.ICPLink),
		CustomTexts: sanitizeFooterCustomTexts(input.CustomTexts),
	}
}

func sanitizeFooterCustomTexts(lines []string) []string {
	if len(lines) == 0 {
		return []string{}
	}

	result := make([]string, 0, min(len(lines), maxFooterCustomTexts))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
		if len(result) >= maxFooterCustomTexts {
			break
		}
	}
	return result
}
