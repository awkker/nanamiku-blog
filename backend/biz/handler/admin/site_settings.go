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
	svc    *service.SiteSettingsService
	modSvc *service.ModerationService
}

func NewSiteSettingsAdminHandler(svc *service.SiteSettingsService, modSvc *service.ModerationService) *SiteSettingsAdminHandler {
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
