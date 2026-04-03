package public

import (
	"context"
	"errors"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type FriendsHandler struct {
	svc    *service.FriendsService
	modSvc *service.ModerationService
}

func NewFriendsHandler(svc *service.FriendsService, modSvc *service.ModerationService) *FriendsHandler {
	return &FriendsHandler{svc: svc, modSvc: modSvc}
}

func (h *FriendsHandler) List(ctx context.Context, c *app.RequestContext) {
	items, err := h.svc.ListApproved(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list friend links"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(items))
}

type createFriendApplicationReq struct {
	SiteName     string `json:"site_name"`
	SiteURL      string `json:"site_url"`
	AvatarURL    string `json:"avatar_url"`
	Description  string `json:"description"`
	ContactEmail string `json:"contact_email"`
	ContactNote  string `json:"contact_note"`
}

func (h *FriendsHandler) CreateApplication(ctx context.Context, c *app.RequestContext) {
	var req createFriendApplicationReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	if h.modSvc != nil {
		word, err := h.modSvc.FindSensitiveWord(ctx, req.SiteName, req.Description, req.ContactNote)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "sensitive-word check failed"))
			return
		}
		if word != "" {
			c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBlocked, "friend application contains blocked keyword"))
			return
		}
	}

	result, err := h.svc.SubmitApplication(ctx, service.FriendApplicationInput{
		SiteName:     req.SiteName,
		SiteURL:      req.SiteURL,
		AvatarURL:    req.AvatarURL,
		Description:  req.Description,
		ContactEmail: req.ContactEmail,
		ContactNote:  req.ContactNote,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidFriendInput):
			c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, err.Error()))
		case errors.Is(err, service.ErrFriendApplicationDuplicate):
			c.JSON(consts.StatusConflict, dto.Err(errcode.ErrConflict, "friend application already exists"))
		default:
			c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to create friend application"))
		}
		return
	}

	c.JSON(consts.StatusCreated, dto.OK(result))
}
