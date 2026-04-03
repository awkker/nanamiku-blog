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
	svc *service.SiteSettingsService
}

func NewSiteSettingsHandler(svc *service.SiteSettingsService) *SiteSettingsHandler {
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
