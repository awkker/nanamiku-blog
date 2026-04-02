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
