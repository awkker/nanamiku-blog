package admin

import (
	"context"
	"errors"
	"strings"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authSvc      *service.AuthService
	cookieConfig AuthCookieConfig
}

func NewAuthHandler(authSvc *service.AuthService, cookieConfig AuthCookieConfig) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, cookieConfig: cookieConfig}
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req loginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request body"))
		return
	}

	identifier := strings.TrimSpace(req.Identifier)
	if identifier == "" {
		identifier = strings.TrimSpace(req.Username)
	}
	if identifier == "" {
		identifier = strings.TrimSpace(req.Email)
	}
	if identifier == "" || strings.TrimSpace(req.Password) == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "identifier and password required"))
		return
	}

	pair, err := h.authSvc.Login(ctx, identifier, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrInvalidCredentials, "invalid username or password"))
			return
		}
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "login failed"))
		return
	}

	h.cookieConfig.writeSessionCookies(c, h.authSvc, pair)
	c.JSON(consts.StatusOK, dto.OK(map[string]interface{}{
		"expires_at": pair.ExpiresAt,
	}))
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" vd:"len($)>0"`
}

func (h *AuthHandler) Refresh(ctx context.Context, c *app.RequestContext) {
	var req refreshRequest
	_ = c.BindJSON(&req)

	refreshToken := resolveRefreshToken(req.RefreshToken, string(c.Cookie(service.RefreshTokenCookieName)))
	if refreshToken == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "refresh_token required"))
		return
	}

	pair, err := h.authSvc.RefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			h.cookieConfig.clearSessionCookies(c)
			c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrTokenExpired, "refresh token expired"))
			return
		}
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "refresh failed"))
		return
	}

	h.cookieConfig.writeSessionCookies(c, h.authSvc, pair)
	c.JSON(consts.StatusOK, dto.OK(map[string]interface{}{
		"expires_at": pair.ExpiresAt,
	}))
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Logout(ctx context.Context, c *app.RequestContext) {
	var req logoutRequest
	_ = c.BindJSON(&req)

	refreshToken := resolveRefreshToken(req.RefreshToken, string(c.Cookie(service.RefreshTokenCookieName)))

	_ = h.authSvc.Logout(ctx, refreshToken)
	h.cookieConfig.clearSessionCookies(c)
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *AuthHandler) PublicProfile(ctx context.Context, c *app.RequestContext) {
	profile, err := h.authSvc.GetPublicProfile(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get public profile"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(map[string]string{
		"username":     profile.Username,
		"display_name": profile.DisplayName,
		"avatar_url":   profile.AvatarURL,
	}))
}

func (h *AuthHandler) Me(ctx context.Context, c *app.RequestContext) {
	adminIDVal, exists := c.Get("admin_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	adminID, ok := adminIDVal.(uuid.UUID)
	if !ok {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	admin, err := h.authSvc.GetAdminInfo(ctx, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get admin info"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(map[string]interface{}{
		"id":            admin.ID,
		"username":      admin.Username,
		"display_name":  admin.DisplayName,
		"avatar_url":    admin.AvatarUrl,
		"email":         admin.Email,
		"role":          admin.Role,
		"last_login_at": admin.LastLoginAt,
		"created_at":    admin.CreatedAt,
	}))
}

type updateMeRequest struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type updateAccountRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) UpdateMe(ctx context.Context, c *app.RequestContext) {
	adminIDVal, exists := c.Get("admin_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	adminID, ok := adminIDVal.(uuid.UUID)
	if !ok {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	var req updateMeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	admin, err := h.authSvc.UpdateProfile(ctx, adminID, req.DisplayName, req.AvatarURL)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to update admin profile"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(map[string]interface{}{
		"id":            admin.ID,
		"username":      admin.Username,
		"display_name":  admin.DisplayName,
		"avatar_url":    admin.AvatarUrl,
		"email":         admin.Email,
		"role":          admin.Role,
		"last_login_at": admin.LastLoginAt,
		"created_at":    admin.CreatedAt,
	}))
}

func (h *AuthHandler) UpdateAccount(ctx context.Context, c *app.RequestContext) {
	adminIDVal, exists := c.Get("admin_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	adminID, ok := adminIDVal.(uuid.UUID)
	if !ok {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	var req updateAccountRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	admin, sessionRevoked, err := h.authSvc.UpdateAccount(ctx, adminID, service.UpdateAccountInput{
		Username:    req.Username,
		Email:       req.Email,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidAccount) || errors.Is(err, service.ErrWeakPassword) {
			c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, err.Error()))
			return
		}
		if errors.Is(err, service.ErrAccountConflict) {
			c.JSON(consts.StatusConflict, dto.Err(errcode.ErrConflict, "username or email already exists"))
			return
		}
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to update account"))
		return
	}

	if sessionRevoked {
		h.cookieConfig.clearSessionCookies(c)
	}

	c.JSON(consts.StatusOK, dto.OK(map[string]interface{}{
		"id":              admin.ID,
		"username":        admin.Username,
		"display_name":    admin.DisplayName,
		"avatar_url":      admin.AvatarUrl,
		"email":           admin.Email,
		"role":            admin.Role,
		"last_login_at":   admin.LastLoginAt,
		"created_at":      admin.CreatedAt,
		"session_revoked": sessionRevoked,
	}))
}

func resolveRefreshToken(requestToken, cookieToken string) string {
	token := strings.TrimSpace(requestToken)
	if token != "" {
		return token
	}

	return strings.TrimSpace(cookieToken)
}
