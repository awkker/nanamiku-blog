package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	siteSettingsFooterKey      = "footer"
	siteSettingsSiteProfileKey = "site_profile"
	maxFooterCustomTexts       = 8
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

type SiteProfileSettings struct {
	BrandText          string `json:"brand_text"`
	LogoAlt            string `json:"logo_alt"`
	SiteTitle          string `json:"site_title"`
	SiteURL            string `json:"site_url"`
	DefaultDescription string `json:"default_description"`
	DefaultSocialImage string `json:"default_social_image"`
}

type siteProfileSettingsDoc struct {
	BrandText          string `json:"brand_text"`
	LogoAlt            string `json:"logo_alt"`
	SiteTitle          string `json:"site_title"`
	SiteURL            string `json:"site_url"`
	DefaultDescription string `json:"default_description"`
	DefaultSocialImage string `json:"default_social_image"`
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

func (s *SiteSettingsService) GetSiteProfileSettings(ctx context.Context) (*SiteProfileSettings, error) {
	row, err := s.q.GetSiteSetting(ctx, siteSettingsSiteProfileKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get site profile settings: %w", err)
	}

	settings, err := decodeSiteProfileSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode site profile settings: %w", err)
	}

	return settings, nil
}

func (s *SiteSettingsService) SaveSiteProfileSettings(ctx context.Context, input SiteProfileSettings, adminID uuid.UUID) (*SiteProfileSettings, error) {
	normalized := normalizeSiteProfileSettings(input)
	payload, err := json.Marshal(siteProfileSettingsDoc{
		BrandText:          normalized.BrandText,
		LogoAlt:            normalized.LogoAlt,
		SiteTitle:          normalized.SiteTitle,
		SiteURL:            normalized.SiteURL,
		DefaultDescription: normalized.DefaultDescription,
		DefaultSocialImage: normalized.DefaultSocialImage,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal site profile settings: %w", err)
	}

	row, err := s.q.UpsertSiteSetting(ctx, query.UpsertSiteSettingParams{
		Key:       siteSettingsSiteProfileKey,
		Value:     payload,
		UpdatedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
	if err != nil {
		return nil, fmt.Errorf("save site profile settings: %w", err)
	}

	settings, err := decodeSiteProfileSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode saved site profile settings: %w", err)
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

func decodeSiteProfileSettings(raw json.RawMessage) (*SiteProfileSettings, error) {
	if len(raw) == 0 {
		return &SiteProfileSettings{}, nil
	}

	var doc siteProfileSettingsDoc
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	normalized := normalizeSiteProfileSettings(SiteProfileSettings{
		BrandText:          doc.BrandText,
		LogoAlt:            doc.LogoAlt,
		SiteTitle:          doc.SiteTitle,
		SiteURL:            doc.SiteURL,
		DefaultDescription: doc.DefaultDescription,
		DefaultSocialImage: doc.DefaultSocialImage,
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

func normalizeSiteProfileSettings(input SiteProfileSettings) SiteProfileSettings {
	brandText := strings.TrimSpace(input.BrandText)
	siteTitle := strings.TrimSpace(input.SiteTitle)
	if siteTitle == "" {
		siteTitle = brandText
	}

	logoAlt := strings.TrimSpace(input.LogoAlt)
	if logoAlt == "" && brandText != "" {
		logoAlt = brandText + " logo"
	}

	return SiteProfileSettings{
		BrandText:          brandText,
		LogoAlt:            logoAlt,
		SiteTitle:          siteTitle,
		SiteURL:            normalizeSiteURL(input.SiteURL),
		DefaultDescription: strings.TrimSpace(input.DefaultDescription),
		DefaultSocialImage: normalizeSiteAssetURL(input.DefaultSocialImage),
	}
}

func normalizeSiteURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "//") {
		trimmed = "https:" + trimmed
	} else if !strings.Contains(trimmed, "://") {
		trimmed = "https://" + strings.TrimLeft(trimmed, "/")
	}

	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}

	parsed.Path = strings.TrimRight(parsed.Path, "/")
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return strings.TrimRight(parsed.String(), "/")
}

func normalizeSiteAssetURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "/") {
		return trimmed
	}

	return normalizeSiteURL(trimmed)
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
