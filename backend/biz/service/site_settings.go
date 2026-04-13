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
	siteSettingsFooterKey       = "footer"
	siteSettingsSiteProfileKey  = "site_profile"
	siteSettingsHomeHeroKey     = "home_hero"
	siteSettingsHomeAssetsKey   = "home_assets"
	siteSettingsAuthorKey       = "author_profile"
	siteSettingsIntegrationsKey = "site_integrations"
	maxFooterCustomTexts        = 8
	maxHomeAssetImages          = 8
	maxAuthorSkills             = 8
	maxAuthorNowItems           = 8
	maxAuthorSocialLinks        = 6
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

type HomeHeroSettings struct {
	HeroTitle    string `json:"hero_title"`
	HeroSubtitle string `json:"hero_subtitle"`
}

type homeHeroSettingsDoc struct {
	HeroTitle    string `json:"hero_title"`
	HeroSubtitle string `json:"hero_subtitle"`
}

type HomeAssetsSettings struct {
	HeroImages []string `json:"hero_images"`
}

type homeAssetsSettingsDoc struct {
	HeroImages []string `json:"hero_images"`
}

type AuthorSocialLink struct {
	Label   string `json:"label"`
	Href    string `json:"href"`
	IconKey string `json:"icon_key"`
}

type AuthorProfileSettings struct {
	DisplayName      string             `json:"display_name"`
	AvatarURL        string             `json:"avatar_url"`
	Role             string             `json:"role"`
	Bio              string             `json:"bio"`
	AboutDescription string             `json:"about_description"`
	Location         string             `json:"location"`
	Since            string             `json:"since"`
	Skills           []string           `json:"skills"`
	NowItems         []string           `json:"now_items"`
	Quote            string             `json:"quote"`
	ContactEmail     string             `json:"contact_email"`
	SocialLinks      []AuthorSocialLink `json:"social_links"`
}

type authorProfileSettingsDoc struct {
	DisplayName      string             `json:"display_name"`
	AvatarURL        string             `json:"avatar_url"`
	Role             string             `json:"role"`
	Bio              string             `json:"bio"`
	AboutDescription string             `json:"about_description"`
	Location         string             `json:"location"`
	Since            string             `json:"since"`
	Skills           []string           `json:"skills"`
	NowItems         []string           `json:"now_items"`
	Quote            string             `json:"quote"`
	ContactEmail     string             `json:"contact_email"`
	SocialLinks      []AuthorSocialLink `json:"social_links"`
}

type SiteIntegrationsSettings struct {
	GitHubUsername  string `json:"github_username"`
	WeatherLocation string `json:"weather_location"`
	ShowWeather     bool   `json:"show_weather"`
	ShowMusic       bool   `json:"show_music"`
	ShowClock       bool   `json:"show_clock"`
}

type siteIntegrationsSettingsDoc struct {
	GitHubUsername  string `json:"github_username"`
	WeatherLocation string `json:"weather_location"`
	ShowWeather     bool   `json:"show_weather"`
	ShowMusic       bool   `json:"show_music"`
	ShowClock       bool   `json:"show_clock"`
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

func (s *SiteSettingsService) GetHomeHeroSettings(ctx context.Context) (*HomeHeroSettings, error) {
	row, err := s.q.GetSiteSetting(ctx, siteSettingsHomeHeroKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get home hero settings: %w", err)
	}

	settings, err := decodeHomeHeroSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode home hero settings: %w", err)
	}

	return settings, nil
}

func (s *SiteSettingsService) SaveHomeHeroSettings(ctx context.Context, input HomeHeroSettings, adminID uuid.UUID) (*HomeHeroSettings, error) {
	normalized := normalizeHomeHeroSettings(input)
	payload, err := json.Marshal(homeHeroSettingsDoc{
		HeroTitle:    normalized.HeroTitle,
		HeroSubtitle: normalized.HeroSubtitle,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal home hero settings: %w", err)
	}

	row, err := s.q.UpsertSiteSetting(ctx, query.UpsertSiteSettingParams{
		Key:       siteSettingsHomeHeroKey,
		Value:     payload,
		UpdatedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
	if err != nil {
		return nil, fmt.Errorf("save home hero settings: %w", err)
	}

	settings, err := decodeHomeHeroSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode saved home hero settings: %w", err)
	}

	return settings, nil
}

func (s *SiteSettingsService) GetHomeAssetsSettings(ctx context.Context) (*HomeAssetsSettings, error) {
	row, err := s.q.GetSiteSetting(ctx, siteSettingsHomeAssetsKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get home assets settings: %w", err)
	}

	settings, err := decodeHomeAssetsSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode home assets settings: %w", err)
	}

	return settings, nil
}

func (s *SiteSettingsService) SaveHomeAssetsSettings(ctx context.Context, input HomeAssetsSettings, adminID uuid.UUID) (*HomeAssetsSettings, error) {
	normalized := normalizeHomeAssetsSettings(input)
	payload, err := json.Marshal(homeAssetsSettingsDoc{
		HeroImages: normalized.HeroImages,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal home assets settings: %w", err)
	}

	row, err := s.q.UpsertSiteSetting(ctx, query.UpsertSiteSettingParams{
		Key:       siteSettingsHomeAssetsKey,
		Value:     payload,
		UpdatedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
	if err != nil {
		return nil, fmt.Errorf("save home assets settings: %w", err)
	}

	settings, err := decodeHomeAssetsSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode saved home assets settings: %w", err)
	}

	return settings, nil
}

func (s *SiteSettingsService) GetAuthorProfileSettings(ctx context.Context) (*AuthorProfileSettings, error) {
	row, err := s.q.GetSiteSetting(ctx, siteSettingsAuthorKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get author profile settings: %w", err)
	}

	settings, err := decodeAuthorProfileSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode author profile settings: %w", err)
	}

	return settings, nil
}

func (s *SiteSettingsService) SaveAuthorProfileSettings(ctx context.Context, input AuthorProfileSettings, adminID uuid.UUID) (*AuthorProfileSettings, error) {
	normalized := normalizeAuthorProfileSettings(input)
	payload, err := json.Marshal(authorProfileSettingsDoc{
		DisplayName:      normalized.DisplayName,
		AvatarURL:        normalized.AvatarURL,
		Role:             normalized.Role,
		Bio:              normalized.Bio,
		AboutDescription: normalized.AboutDescription,
		Location:         normalized.Location,
		Since:            normalized.Since,
		Skills:           normalized.Skills,
		NowItems:         normalized.NowItems,
		Quote:            normalized.Quote,
		ContactEmail:     normalized.ContactEmail,
		SocialLinks:      normalized.SocialLinks,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal author profile settings: %w", err)
	}

	row, err := s.q.UpsertSiteSetting(ctx, query.UpsertSiteSettingParams{
		Key:       siteSettingsAuthorKey,
		Value:     payload,
		UpdatedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
	if err != nil {
		return nil, fmt.Errorf("save author profile settings: %w", err)
	}

	settings, err := decodeAuthorProfileSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode saved author profile settings: %w", err)
	}

	return settings, nil
}

func (s *SiteSettingsService) GetSiteIntegrationsSettings(ctx context.Context) (*SiteIntegrationsSettings, error) {
	row, err := s.q.GetSiteSetting(ctx, siteSettingsIntegrationsKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get site integrations settings: %w", err)
	}

	settings, err := decodeSiteIntegrationsSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode site integrations settings: %w", err)
	}

	return settings, nil
}

func (s *SiteSettingsService) SaveSiteIntegrationsSettings(ctx context.Context, input SiteIntegrationsSettings, adminID uuid.UUID) (*SiteIntegrationsSettings, error) {
	normalized := normalizeSiteIntegrationsSettings(input)
	payload, err := json.Marshal(siteIntegrationsSettingsDoc{
		GitHubUsername:  normalized.GitHubUsername,
		WeatherLocation: normalized.WeatherLocation,
		ShowWeather:     normalized.ShowWeather,
		ShowMusic:       normalized.ShowMusic,
		ShowClock:       normalized.ShowClock,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal site integrations settings: %w", err)
	}

	row, err := s.q.UpsertSiteSetting(ctx, query.UpsertSiteSettingParams{
		Key:       siteSettingsIntegrationsKey,
		Value:     payload,
		UpdatedBy: pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
	})
	if err != nil {
		return nil, fmt.Errorf("save site integrations settings: %w", err)
	}

	settings, err := decodeSiteIntegrationsSettings(row.Value)
	if err != nil {
		return nil, fmt.Errorf("decode saved site integrations settings: %w", err)
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

func decodeHomeHeroSettings(raw json.RawMessage) (*HomeHeroSettings, error) {
	if len(raw) == 0 {
		return &HomeHeroSettings{}, nil
	}

	var doc homeHeroSettingsDoc
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	normalized := normalizeHomeHeroSettings(HomeHeroSettings{
		HeroTitle:    doc.HeroTitle,
		HeroSubtitle: doc.HeroSubtitle,
	})

	return &normalized, nil
}

func decodeHomeAssetsSettings(raw json.RawMessage) (*HomeAssetsSettings, error) {
	if len(raw) == 0 {
		return &HomeAssetsSettings{HeroImages: []string{}}, nil
	}

	var doc homeAssetsSettingsDoc
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	normalized := normalizeHomeAssetsSettings(HomeAssetsSettings{
		HeroImages: doc.HeroImages,
	})

	return &normalized, nil
}

func decodeAuthorProfileSettings(raw json.RawMessage) (*AuthorProfileSettings, error) {
	if len(raw) == 0 {
		return &AuthorProfileSettings{
			Skills:      []string{},
			NowItems:    []string{},
			SocialLinks: []AuthorSocialLink{},
		}, nil
	}

	var doc authorProfileSettingsDoc
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	normalized := normalizeAuthorProfileSettings(AuthorProfileSettings{
		DisplayName:      doc.DisplayName,
		AvatarURL:        doc.AvatarURL,
		Role:             doc.Role,
		Bio:              doc.Bio,
		AboutDescription: doc.AboutDescription,
		Location:         doc.Location,
		Since:            doc.Since,
		Skills:           doc.Skills,
		NowItems:         doc.NowItems,
		Quote:            doc.Quote,
		ContactEmail:     doc.ContactEmail,
		SocialLinks:      doc.SocialLinks,
	})

	return &normalized, nil
}

func decodeSiteIntegrationsSettings(raw json.RawMessage) (*SiteIntegrationsSettings, error) {
	if len(raw) == 0 {
		return &SiteIntegrationsSettings{}, nil
	}

	var doc siteIntegrationsSettingsDoc
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	normalized := normalizeSiteIntegrationsSettings(SiteIntegrationsSettings{
		GitHubUsername:  doc.GitHubUsername,
		WeatherLocation: doc.WeatherLocation,
		ShowWeather:     doc.ShowWeather,
		ShowMusic:       doc.ShowMusic,
		ShowClock:       doc.ShowClock,
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

func normalizeHomeHeroSettings(input HomeHeroSettings) HomeHeroSettings {
	return HomeHeroSettings{
		HeroTitle:    strings.TrimSpace(input.HeroTitle),
		HeroSubtitle: strings.TrimSpace(input.HeroSubtitle),
	}
}

func normalizeHomeAssetsSettings(input HomeAssetsSettings) HomeAssetsSettings {
	return HomeAssetsSettings{
		HeroImages: sanitizeHomeAssetImages(input.HeroImages),
	}
}

func normalizeAuthorProfileSettings(input AuthorProfileSettings) AuthorProfileSettings {
	return AuthorProfileSettings{
		DisplayName:      strings.TrimSpace(input.DisplayName),
		AvatarURL:        normalizeAuthorAvatarURL(input.AvatarURL),
		Role:             strings.TrimSpace(input.Role),
		Bio:              strings.TrimSpace(input.Bio),
		AboutDescription: strings.TrimSpace(input.AboutDescription),
		Location:         strings.TrimSpace(input.Location),
		Since:            strings.TrimSpace(input.Since),
		Skills:           sanitizeStringItems(input.Skills, maxAuthorSkills),
		NowItems:         sanitizeStringItems(input.NowItems, maxAuthorNowItems),
		Quote:            strings.TrimSpace(input.Quote),
		ContactEmail:     sanitizeContactEmail(input.ContactEmail),
		SocialLinks:      sanitizeAuthorSocialLinks(input.SocialLinks),
	}
}

func normalizeAuthorAvatarURL(raw string) string {
	return normalizeSiteAssetURL(raw)
}

func normalizeSiteIntegrationsSettings(input SiteIntegrationsSettings) SiteIntegrationsSettings {
	return SiteIntegrationsSettings{
		GitHubUsername:  strings.TrimSpace(input.GitHubUsername),
		WeatherLocation: strings.TrimSpace(input.WeatherLocation),
		ShowWeather:     input.ShowWeather,
		ShowMusic:       input.ShowMusic,
		ShowClock:       input.ShowClock,
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

	if strings.HasPrefix(trimmed, "./") || strings.HasPrefix(trimmed, "../") || looksLikeLocalAssetPath(trimmed) {
		return normalizeRelativeAssetPath(trimmed)
	}

	normalized := normalizeSiteURL(trimmed)
	if normalized == "" {
		return ""
	}

	if (strings.HasPrefix(normalized, "https://picture/") || strings.HasPrefix(normalized, "http://picture/")) && hasImageAssetSuffix(normalized) {
		return strings.Replace(strings.Replace(normalized, "https://picture/", "/picture/", 1), "http://picture/", "/picture/", 1)
	}

	return normalized
}

func normalizeRelativeAssetPath(raw string) string {
	trimmed := strings.TrimSpace(raw)
	for {
		switch {
		case strings.HasPrefix(trimmed, "./"):
			trimmed = strings.TrimPrefix(trimmed, "./")
		case strings.HasPrefix(trimmed, "../"):
			trimmed = strings.TrimPrefix(trimmed, "../")
		default:
			return "/" + strings.TrimLeft(trimmed, "/")
		}
	}
}

func looksLikeLocalAssetPath(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || strings.Contains(trimmed, "://") || !strings.Contains(trimmed, "/") || !hasImageAssetSuffix(trimmed) {
		return false
	}

	head, _, _ := strings.Cut(trimmed, "/")
	return head != "" && !strings.Contains(head, ".")
}

func hasImageAssetSuffix(raw string) bool {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	for _, suffix := range []string{".avif", ".gif", ".jpg", ".jpeg", ".png", ".svg", ".webp"} {
		if strings.HasSuffix(normalized, suffix) {
			return true
		}
	}
	return false
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

func sanitizeHomeAssetImages(images []string) []string {
	if len(images) == 0 {
		return []string{}
	}

	result := make([]string, 0, min(len(images), maxHomeAssetImages))
	seen := make(map[string]struct{}, len(images))

	for _, image := range images {
		normalized := normalizeSiteAssetURL(image)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
		if len(result) >= maxHomeAssetImages {
			break
		}
	}

	return result
}

func sanitizeStringItems(items []string, limit int) []string {
	if len(items) == 0 {
		return []string{}
	}

	result := make([]string, 0, min(len(items), limit))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
		if len(result) >= limit {
			break
		}
	}

	return result
}

func sanitizeContactEmail(email string) string {
	trimmed := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(email), "mailto:"))
	return trimmed
}

func sanitizeAuthorSocialLinks(items []AuthorSocialLink) []AuthorSocialLink {
	if len(items) == 0 {
		return []AuthorSocialLink{}
	}

	result := make([]AuthorSocialLink, 0, min(len(items), maxAuthorSocialLinks))
	for _, item := range items {
		label := strings.TrimSpace(item.Label)
		href := normalizePublicLink(item.Href)
		if label == "" || href == "" {
			continue
		}
		result = append(result, AuthorSocialLink{
			Label:   label,
			Href:    href,
			IconKey: strings.TrimSpace(item.IconKey),
		})
		if len(result) >= maxAuthorSocialLinks {
			break
		}
	}

	return result
}

func normalizePublicLink(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "/") ||
		strings.HasPrefix(trimmed, "http://") ||
		strings.HasPrefix(trimmed, "https://") ||
		strings.HasPrefix(trimmed, "mailto:") ||
		strings.HasPrefix(trimmed, "tel:") {
		return trimmed
	}

	return normalizeSiteURL(trimmed)
}
