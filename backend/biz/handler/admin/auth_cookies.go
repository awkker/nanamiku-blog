package admin

import (
	"strings"

	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type AuthCookieConfig struct {
	Domain   string
	Secure   bool
	SameSite string
}

func (cfg AuthCookieConfig) sameSiteMode() protocol.CookieSameSite {
	switch strings.ToLower(strings.TrimSpace(cfg.SameSite)) {
	case "strict":
		return protocol.CookieSameSiteStrictMode
	case "none":
		return protocol.CookieSameSiteNoneMode
	default:
		return protocol.CookieSameSiteLaxMode
	}
}

func (cfg AuthCookieConfig) writeSessionCookies(c *app.RequestContext, authSvc *service.AuthService, pair *service.TokenPair) {
	c.SetCookie(
		service.AccessTokenCookieName,
		pair.AccessToken,
		int(authSvc.AccessTTL().Seconds()),
		"/",
		cfg.Domain,
		cfg.sameSiteMode(),
		cfg.Secure,
		true,
	)
	c.SetCookie(
		service.RefreshTokenCookieName,
		pair.RefreshToken,
		int(authSvc.RefreshTTL().Seconds()),
		"/",
		cfg.Domain,
		cfg.sameSiteMode(),
		cfg.Secure,
		true,
	)
}

func (cfg AuthCookieConfig) clearSessionCookies(c *app.RequestContext) {
	c.SetCookie(
		service.AccessTokenCookieName,
		"",
		-1,
		"/",
		cfg.Domain,
		cfg.sameSiteMode(),
		cfg.Secure,
		true,
	)
	c.SetCookie(
		service.RefreshTokenCookieName,
		"",
		-1,
		"/",
		cfg.Domain,
		cfg.sameSiteMode(),
		cfg.Secure,
		true,
	)
}
