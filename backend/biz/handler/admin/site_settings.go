package admin

import (
	"context"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type SiteSettingsAdminHandler struct {
	svc    siteSettingsAdminService
	modSvc siteSettingsAuditLogger
}

type siteSettingsAdminService interface {
	SaveFooterSettings(ctx context.Context, input service.FooterSettings, adminID uuid.UUID) (*service.FooterSettings, error)
	SaveSiteProfileSettings(ctx context.Context, input service.SiteProfileSettings, adminID uuid.UUID) (*service.SiteProfileSettings, error)
	SaveHomeHeroSettings(ctx context.Context, input service.HomeHeroSettings, adminID uuid.UUID) (*service.HomeHeroSettings, error)
	SaveHomeAssetsSettings(ctx context.Context, input service.HomeAssetsSettings, adminID uuid.UUID) (*service.HomeAssetsSettings, error)
	SaveBlogIndexSettings(ctx context.Context, input service.BlogIndexSettings, adminID uuid.UUID) (*service.BlogIndexSettings, error)
	SaveAuthorProfileSettings(ctx context.Context, input service.AuthorProfileSettings, adminID uuid.UUID) (*service.AuthorProfileSettings, error)
	SaveSiteIntegrationsSettings(ctx context.Context, input service.SiteIntegrationsSettings, adminID uuid.UUID) (*service.SiteIntegrationsSettings, error)
}

type siteSettingsAuditLogger interface {
	LogAudit(ctx context.Context, adminID uuid.UUID, action, targetType, targetID string, detail interface{}, ip string) error
}

func NewSiteSettingsAdminHandler(svc siteSettingsAdminService, modSvc siteSettingsAuditLogger) *SiteSettingsAdminHandler {
	return &SiteSettingsAdminHandler{svc: svc, modSvc: modSvc}
}

type updateFooterSettingsReq struct {
	ICPText     string   `json:"icp_text"`
	ICPLink     string   `json:"icp_link"`
	CustomTexts []string `json:"custom_texts"`
}

type updateSiteProfileSettingsReq struct {
	BrandText          string `json:"brand_text"`
	LogoAlt            string `json:"logo_alt"`
	SiteTitle          string `json:"site_title"`
	SiteURL            string `json:"site_url"`
	DefaultDescription string `json:"default_description"`
	DefaultSocialImage string `json:"default_social_image"`
}

type updateHomeHeroSettingsReq struct {
	HeroTitle    string `json:"hero_title"`
	HeroSubtitle string `json:"hero_subtitle"`
}

type updateHomeAssetsSettingsReq struct {
	HeroImages []string `json:"hero_images"`
}

type updateBlogIndexActionReq struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

type updateBlogIndexStatReq struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type updateBlogIndexFocusCardReq struct {
	Badge       string `json:"badge"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Footnote    string `json:"footnote"`
}

type updateBlogIndexScrollCueReq struct {
	Label     string `json:"label"`
	AriaLabel string `json:"aria_label"`
}

type updateBlogIndexSettingsReq struct {
	HeroBadge       string                      `json:"hero_badge"`
	HeroTitle       string                      `json:"hero_title"`
	HeroDescription string                      `json:"hero_description"`
	HeroActions     []updateBlogIndexActionReq  `json:"hero_actions"`
	QuickStats      []updateBlogIndexStatReq    `json:"quick_stats"`
	FocusCard       updateBlogIndexFocusCardReq `json:"focus_card"`
	ScrollCue       updateBlogIndexScrollCueReq `json:"scroll_cue"`
}

type updateAuthorSocialLinkReq struct {
	Label   string `json:"label"`
	Href    string `json:"href"`
	IconKey string `json:"icon_key"`
}

type updateAuthorProfileSettingsReq struct {
	DisplayName      string                      `json:"display_name"`
	AvatarURL        string                      `json:"avatar_url"`
	Role             string                      `json:"role"`
	Bio              string                      `json:"bio"`
	AboutDescription string                      `json:"about_description"`
	Location         string                      `json:"location"`
	Since            string                      `json:"since"`
	Skills           []string                    `json:"skills"`
	NowItems         []string                    `json:"now_items"`
	Quote            string                      `json:"quote"`
	ContactEmail     string                      `json:"contact_email"`
	SocialLinks      []updateAuthorSocialLinkReq `json:"social_links"`
}

type updateSiteIntegrationsSettingsReq struct {
	GitHubUsername  string `json:"github_username"`
	WeatherLocation string `json:"weather_location"`
	ShowWeather     bool   `json:"show_weather"`
	ShowMusic       bool   `json:"show_music"`
	ShowClock       bool   `json:"show_clock"`
}

func (h *SiteSettingsAdminHandler) UpdateFooter(ctx context.Context, c *app.RequestContext) {
	adminID := getAdminID(c)
	if adminID == uuid.Nil {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	var req updateFooterSettingsReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	settings, err := h.svc.SaveFooterSettings(ctx, service.FooterSettings{
		ICPText:     req.ICPText,
		ICPLink:     req.ICPLink,
		CustomTexts: req.CustomTexts,
	}, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to save footer settings"))
		return
	}

	_ = h.modSvc.LogAudit(ctx, adminID, "update", "site_setting", "footer", map[string]interface{}{
		"icp_text":     settings.ICPText,
		"icp_link":     settings.ICPLink,
		"custom_count": len(settings.CustomTexts),
	}, getClientIP(c))

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsAdminHandler) UpdateSiteProfile(ctx context.Context, c *app.RequestContext) {
	adminID := getAdminID(c)
	if adminID == uuid.Nil {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	var req updateSiteProfileSettingsReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	settings, err := h.svc.SaveSiteProfileSettings(ctx, service.SiteProfileSettings{
		BrandText:          req.BrandText,
		LogoAlt:            req.LogoAlt,
		SiteTitle:          req.SiteTitle,
		SiteURL:            req.SiteURL,
		DefaultDescription: req.DefaultDescription,
		DefaultSocialImage: req.DefaultSocialImage,
	}, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to save site profile settings"))
		return
	}

	_ = h.modSvc.LogAudit(ctx, adminID, "update", "site_setting", "site_profile", map[string]interface{}{
		"brand_text":           settings.BrandText,
		"site_title":           settings.SiteTitle,
		"site_url":             settings.SiteURL,
		"default_description":  settings.DefaultDescription,
		"default_social_image": settings.DefaultSocialImage,
	}, getClientIP(c))

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsAdminHandler) UpdateHomeHero(ctx context.Context, c *app.RequestContext) {
	adminID := getAdminID(c)
	if adminID == uuid.Nil {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	var req updateHomeHeroSettingsReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	settings, err := h.svc.SaveHomeHeroSettings(ctx, service.HomeHeroSettings{
		HeroTitle:    req.HeroTitle,
		HeroSubtitle: req.HeroSubtitle,
	}, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to save home hero settings"))
		return
	}

	_ = h.modSvc.LogAudit(ctx, adminID, "update", "site_setting", "home_hero", map[string]interface{}{
		"hero_title":    settings.HeroTitle,
		"hero_subtitle": settings.HeroSubtitle,
	}, getClientIP(c))

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsAdminHandler) UpdateHomeAssets(ctx context.Context, c *app.RequestContext) {
	adminID := getAdminID(c)
	if adminID == uuid.Nil {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	var req updateHomeAssetsSettingsReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	settings, err := h.svc.SaveHomeAssetsSettings(ctx, service.HomeAssetsSettings{
		HeroImages: req.HeroImages,
	}, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to save home assets settings"))
		return
	}

	_ = h.modSvc.LogAudit(ctx, adminID, "update", "site_setting", "home_assets", map[string]interface{}{
		"hero_image_count": len(settings.HeroImages),
		"hero_images":      settings.HeroImages,
	}, getClientIP(c))

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsAdminHandler) UpdateBlogIndex(ctx context.Context, c *app.RequestContext) {
	adminID := getAdminID(c)
	if adminID == uuid.Nil {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	var req updateBlogIndexSettingsReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	heroActions := make([]service.BlogIndexHeroAction, 0, len(req.HeroActions))
	for _, item := range req.HeroActions {
		heroActions = append(heroActions, service.BlogIndexHeroAction{
			Label: item.Label,
			Href:  item.Href,
		})
	}

	quickStats := make([]service.BlogIndexQuickStat, 0, len(req.QuickStats))
	for _, item := range req.QuickStats {
		quickStats = append(quickStats, service.BlogIndexQuickStat{
			Label: item.Label,
			Value: item.Value,
		})
	}

	settings, err := h.svc.SaveBlogIndexSettings(ctx, service.BlogIndexSettings{
		HeroBadge:       req.HeroBadge,
		HeroTitle:       req.HeroTitle,
		HeroDescription: req.HeroDescription,
		HeroActions:     heroActions,
		QuickStats:      quickStats,
		FocusCard: service.BlogIndexFocusCard{
			Badge:       req.FocusCard.Badge,
			Title:       req.FocusCard.Title,
			Description: req.FocusCard.Description,
			Footnote:    req.FocusCard.Footnote,
		},
		ScrollCue: service.BlogIndexScrollCue{
			Label:     req.ScrollCue.Label,
			AriaLabel: req.ScrollCue.AriaLabel,
		},
	}, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to save blog index settings"))
		return
	}

	_ = h.modSvc.LogAudit(ctx, adminID, "update", "site_setting", "blog_index", map[string]interface{}{
		"hero_title":         settings.HeroTitle,
		"hero_actions_count": len(settings.HeroActions),
		"quick_stats_count":  len(settings.QuickStats),
		"focus_title":        settings.FocusCard.Title,
	}, getClientIP(c))

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsAdminHandler) UpdateAuthorProfile(ctx context.Context, c *app.RequestContext) {
	adminID := getAdminID(c)
	if adminID == uuid.Nil {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	var req updateAuthorProfileSettingsReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	socialLinks := make([]service.AuthorSocialLink, 0, len(req.SocialLinks))
	for _, item := range req.SocialLinks {
		socialLinks = append(socialLinks, service.AuthorSocialLink{
			Label:   item.Label,
			Href:    item.Href,
			IconKey: item.IconKey,
		})
	}

	settings, err := h.svc.SaveAuthorProfileSettings(ctx, service.AuthorProfileSettings{
		DisplayName:      req.DisplayName,
		AvatarURL:        req.AvatarURL,
		Role:             req.Role,
		Bio:              req.Bio,
		AboutDescription: req.AboutDescription,
		Location:         req.Location,
		Since:            req.Since,
		Skills:           req.Skills,
		NowItems:         req.NowItems,
		Quote:            req.Quote,
		ContactEmail:     req.ContactEmail,
		SocialLinks:      socialLinks,
	}, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to save author profile settings"))
		return
	}

	_ = h.modSvc.LogAudit(ctx, adminID, "update", "site_setting", "author_profile", map[string]interface{}{
		"display_name":       settings.DisplayName,
		"role":               settings.Role,
		"skills_count":       len(settings.Skills),
		"social_links_count": len(settings.SocialLinks),
	}, getClientIP(c))

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsAdminHandler) UpdateSiteIntegrations(ctx context.Context, c *app.RequestContext) {
	adminID := getAdminID(c)
	if adminID == uuid.Nil {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	var req updateSiteIntegrationsSettingsReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	settings, err := h.svc.SaveSiteIntegrationsSettings(ctx, service.SiteIntegrationsSettings{
		GitHubUsername:  req.GitHubUsername,
		WeatherLocation: req.WeatherLocation,
		ShowWeather:     req.ShowWeather,
		ShowMusic:       req.ShowMusic,
		ShowClock:       req.ShowClock,
	}, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to save site integrations settings"))
		return
	}

	_ = h.modSvc.LogAudit(ctx, adminID, "update", "site_setting", "site_integrations", map[string]interface{}{
		"github_username":  settings.GitHubUsername,
		"weather_location": settings.WeatherLocation,
		"show_weather":     settings.ShowWeather,
		"show_music":       settings.ShowMusic,
		"show_clock":       settings.ShowClock,
	}, getClientIP(c))

	c.JSON(consts.StatusOK, dto.OK(settings))
}
