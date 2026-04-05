package public

import (
	"context"
	"log/slog"
	"strings"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type GitHubHandler struct {
	svc *service.GitHubProfileService
}

func NewGitHubHandler(svc *service.GitHubProfileService) *GitHubHandler {
	return &GitHubHandler{svc: svc}
}

func (h *GitHubHandler) Profile(ctx context.Context, c *app.RequestContext) {
	username := strings.TrimSpace(c.Query("username"))
	if username == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "username required"))
		return
	}

	data, err := h.svc.GetProfile(ctx, username)
	if err != nil {
		slog.Error("failed to fetch github profile", "username", username, "error", err)
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to fetch github profile"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(data))
}
