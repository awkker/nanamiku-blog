package public

import (
	"context"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type SiteSettingsHandler struct {
	svc siteSettingsService
}

type siteSettingsService interface {
	GetFooterSettings(ctx context.Context) (*service.FooterSettings, error)
	GetSiteProfileSettings(ctx context.Context) (*service.SiteProfileSettings, error)
	GetHomeHeroSettings(ctx context.Context) (*service.HomeHeroSettings, error)
	GetHomeAssetsSettings(ctx context.Context) (*service.HomeAssetsSettings, error)
	GetAuthorProfileSettings(ctx context.Context) (*service.AuthorProfileSettings, error)
	GetSiteIntegrationsSettings(ctx context.Context) (*service.SiteIntegrationsSettings, error)
}

func NewSiteSettingsHandler(svc siteSettingsService) *SiteSettingsHandler {
	return &SiteSettingsHandler{svc: svc}
}

func (h *SiteSettingsHandler) GetFooter(ctx context.Context, c *app.RequestContext) {
	settings, err := h.svc.GetFooterSettings(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get footer settings"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsHandler) GetSiteProfile(ctx context.Context, c *app.RequestContext) {
	settings, err := h.svc.GetSiteProfileSettings(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get site profile settings"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsHandler) GetHomeHero(ctx context.Context, c *app.RequestContext) {
	settings, err := h.svc.GetHomeHeroSettings(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get home hero settings"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsHandler) GetHomeAssets(ctx context.Context, c *app.RequestContext) {
	settings, err := h.svc.GetHomeAssetsSettings(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get home assets settings"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsHandler) GetAuthorProfile(ctx context.Context, c *app.RequestContext) {
	settings, err := h.svc.GetAuthorProfileSettings(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get author profile settings"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(settings))
}

func (h *SiteSettingsHandler) GetSiteIntegrations(ctx context.Context, c *app.RequestContext) {
	settings, err := h.svc.GetSiteIntegrationsSettings(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get site integrations settings"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(settings))
}
