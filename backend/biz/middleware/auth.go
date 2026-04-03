package middleware

import (
	"context"
	"strings"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type AdminClaims struct {
	AdminID  uuid.UUID
	Username string
	Role     string
}

type TokenValidator func(tokenStr string) (*AdminClaims, error)

func AdminAuth(validate TokenValidator) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		tokenStr := resolveAccessToken(
			string(c.GetHeader("Authorization")),
			string(c.Cookie(service.AccessTokenCookieName)),
		)
		if tokenStr == "" {
			c.AbortWithStatusJSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "missing authentication credentials"))
			return
		}
		claims, err := validate(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(consts.StatusUnauthorized, dto.Err(errcode.ErrTokenInvalid, "invalid or expired token"))
			return
		}

		c.Set("admin_id", claims.AdminID)
		c.Set("admin_username", claims.Username)
		c.Set("admin_role", claims.Role)
		c.Next(ctx)
	}
}

func resolveAccessToken(authorizationHeader, accessCookie string) string {
	header := strings.TrimSpace(authorizationHeader)
	if strings.HasPrefix(header, "Bearer ") {
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if token != "" {
			return token
		}
	}

	return strings.TrimSpace(accessCookie)
}
