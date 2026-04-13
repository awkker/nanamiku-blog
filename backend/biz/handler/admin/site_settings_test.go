package admin

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type stubSiteSettingsAdminService struct {
	saveFooter           func(context.Context, service.FooterSettings, uuid.UUID) (*service.FooterSettings, error)
	saveSiteProfile      func(context.Context, service.SiteProfileSettings, uuid.UUID) (*service.SiteProfileSettings, error)
	saveHomeHero         func(context.Context, service.HomeHeroSettings, uuid.UUID) (*service.HomeHeroSettings, error)
	saveHomeAssets       func(context.Context, service.HomeAssetsSettings, uuid.UUID) (*service.HomeAssetsSettings, error)
	saveBlogIndex        func(context.Context, service.BlogIndexSettings, uuid.UUID) (*service.BlogIndexSettings, error)
	saveAuthorProfile    func(context.Context, service.AuthorProfileSettings, uuid.UUID) (*service.AuthorProfileSettings, error)
	saveSiteIntegrations func(context.Context, service.SiteIntegrationsSettings, uuid.UUID) (*service.SiteIntegrationsSettings, error)
}

func (s *stubSiteSettingsAdminService) SaveFooterSettings(ctx context.Context, input service.FooterSettings, adminID uuid.UUID) (*service.FooterSettings, error) {
	if s.saveFooter != nil {
		return s.saveFooter(ctx, input, adminID)
	}
	return &service.FooterSettings{}, nil
}

func (s *stubSiteSettingsAdminService) SaveSiteProfileSettings(ctx context.Context, input service.SiteProfileSettings, adminID uuid.UUID) (*service.SiteProfileSettings, error) {
	if s.saveSiteProfile != nil {
		return s.saveSiteProfile(ctx, input, adminID)
	}
	return &service.SiteProfileSettings{}, nil
}

func (s *stubSiteSettingsAdminService) SaveHomeHeroSettings(ctx context.Context, input service.HomeHeroSettings, adminID uuid.UUID) (*service.HomeHeroSettings, error) {
	if s.saveHomeHero != nil {
		return s.saveHomeHero(ctx, input, adminID)
	}
	return &service.HomeHeroSettings{}, nil
}

func (s *stubSiteSettingsAdminService) SaveHomeAssetsSettings(ctx context.Context, input service.HomeAssetsSettings, adminID uuid.UUID) (*service.HomeAssetsSettings, error) {
	if s.saveHomeAssets != nil {
		return s.saveHomeAssets(ctx, input, adminID)
	}
	return &service.HomeAssetsSettings{}, nil
}

func (s *stubSiteSettingsAdminService) SaveBlogIndexSettings(ctx context.Context, input service.BlogIndexSettings, adminID uuid.UUID) (*service.BlogIndexSettings, error) {
	if s.saveBlogIndex != nil {
		return s.saveBlogIndex(ctx, input, adminID)
	}
	return &service.BlogIndexSettings{}, nil
}

func (s *stubSiteSettingsAdminService) SaveAuthorProfileSettings(ctx context.Context, input service.AuthorProfileSettings, adminID uuid.UUID) (*service.AuthorProfileSettings, error) {
	if s.saveAuthorProfile != nil {
		return s.saveAuthorProfile(ctx, input, adminID)
	}
	return &service.AuthorProfileSettings{}, nil
}

func (s *stubSiteSettingsAdminService) SaveSiteIntegrationsSettings(ctx context.Context, input service.SiteIntegrationsSettings, adminID uuid.UUID) (*service.SiteIntegrationsSettings, error) {
	if s.saveSiteIntegrations != nil {
		return s.saveSiteIntegrations(ctx, input, adminID)
	}
	return &service.SiteIntegrationsSettings{}, nil
}

type stubAuditLogger struct {
	calls []auditCall
	err   error
}

type auditCall struct {
	adminID    uuid.UUID
	action     string
	targetType string
	targetID   string
	detail     interface{}
	ip         string
}

func (s *stubAuditLogger) LogAudit(_ context.Context, adminID uuid.UUID, action, targetType, targetID string, detail interface{}, ip string) error {
	s.calls = append(s.calls, auditCall{
		adminID:    adminID,
		action:     action,
		targetType: targetType,
		targetID:   targetID,
		detail:     detail,
		ip:         ip,
	})
	return s.err
}

type adminResponseEnvelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func newAdminTestContext(method, uri, body string) *app.RequestContext {
	ctx := app.NewContext(0)
	ctx.Request.SetMethod(method)
	ctx.Request.SetRequestURI(uri)
	ctx.Request.Header.Set("X-Real-IP", "127.0.0.1")
	if body != "" {
		ctx.Request.Header.Set("Content-Type", consts.MIMEApplicationJSON)
		ctx.Request.SetBodyString(body)
	}
	return ctx
}

func decodeAdminResponse(t *testing.T, ctx *app.RequestContext) adminResponseEnvelope {
	t.Helper()

	var body adminResponseEnvelope
	if err := json.Unmarshal(ctx.Response.Body(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return body
}

func TestSiteSettingsAdminHandlerUpdateSiteProfileUnauthorized(t *testing.T) {
	t.Parallel()

	handler := NewSiteSettingsAdminHandler(&stubSiteSettingsAdminService{}, &stubAuditLogger{})
	ctx := newAdminTestContext(consts.MethodPut, "/api/v1/admin/site-settings/site-profile", `{"brand_text":"Miku"}`)

	handler.UpdateSiteProfile(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusUnauthorized {
		t.Fatalf("unexpected status code: %d", ctx.Response.StatusCode())
	}
}

func TestSiteSettingsAdminHandlerUpdateSiteProfileSuccess(t *testing.T) {
	t.Parallel()

	adminID := uuid.New()
	logger := &stubAuditLogger{}
	var received service.SiteProfileSettings
	var receivedAdminID uuid.UUID

	handler := NewSiteSettingsAdminHandler(&stubSiteSettingsAdminService{
		saveSiteProfile: func(_ context.Context, input service.SiteProfileSettings, inputAdminID uuid.UUID) (*service.SiteProfileSettings, error) {
			received = input
			receivedAdminID = inputAdminID
			return &service.SiteProfileSettings{
				BrandText:          "Miku Blog",
				SiteTitle:          "Miku Blog",
				LogoAlt:            "Miku Blog logo",
				SiteURL:            "https://example.com",
				DefaultDescription: "hello",
				DefaultSocialImage: "/picture/social.png",
			}, nil
		},
	}, logger)

	ctx := newAdminTestContext(
		consts.MethodPut,
		"/api/v1/admin/site-settings/site-profile",
		`{"brand_text":" Miku Blog ","logo_alt":" Miku Blog logo ","site_title":" Miku Blog ","site_url":" https://example.com ","default_description":" hello ","default_social_image":" /picture/social.png "}`,
	)
	ctx.Set("admin_id", adminID)

	handler.UpdateSiteProfile(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusOK {
		t.Fatalf("unexpected status code: %d", ctx.Response.StatusCode())
	}

	if receivedAdminID != adminID {
		t.Fatalf("unexpected admin id: %s", receivedAdminID)
	}
	if received.BrandText != " Miku Blog " {
		t.Fatalf("request payload was not bound as expected: %+v", received)
	}
	if len(logger.calls) != 1 {
		t.Fatalf("expected one audit log call, got %d", len(logger.calls))
	}
	if logger.calls[0].targetID != "site_profile" {
		t.Fatalf("unexpected audit target: %+v", logger.calls[0])
	}

	resp := decodeAdminResponse(t, ctx)
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d", resp.Code)
	}
}

func TestSiteSettingsAdminHandlerUpdateSiteIntegrationsInvalidBody(t *testing.T) {
	t.Parallel()

	handler := NewSiteSettingsAdminHandler(&stubSiteSettingsAdminService{}, &stubAuditLogger{})
	ctx := newAdminTestContext(consts.MethodPut, "/api/v1/admin/site-settings/site-integrations", `{`)
	ctx.Set("admin_id", uuid.New())

	handler.UpdateSiteIntegrations(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", ctx.Response.StatusCode())
	}
}

func TestSiteSettingsAdminHandlerUpdateSiteIntegrationsServiceError(t *testing.T) {
	t.Parallel()

	handler := NewSiteSettingsAdminHandler(&stubSiteSettingsAdminService{
		saveSiteIntegrations: func(context.Context, service.SiteIntegrationsSettings, uuid.UUID) (*service.SiteIntegrationsSettings, error) {
			return nil, errors.New("boom")
		},
	}, &stubAuditLogger{})

	ctx := newAdminTestContext(
		consts.MethodPut,
		"/api/v1/admin/site-settings/site-integrations",
		`{"github_username":"yourname","weather_location":"Tokyo","show_weather":true,"show_music":false,"show_clock":true}`,
	)
	ctx.Set("admin_id", uuid.New())

	handler.UpdateSiteIntegrations(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", ctx.Response.StatusCode())
	}

	resp := decodeAdminResponse(t, ctx)
	if resp.Message != "failed to save site integrations settings" {
		t.Fatalf("unexpected response message: %q", resp.Message)
	}
}

func TestSiteSettingsAdminHandlerUpdateBlogIndexSuccess(t *testing.T) {
	t.Parallel()

	adminID := uuid.New()
	logger := &stubAuditLogger{}
	var received service.BlogIndexSettings

	handler := NewSiteSettingsAdminHandler(&stubSiteSettingsAdminService{
		saveBlogIndex: func(_ context.Context, input service.BlogIndexSettings, inputAdminID uuid.UUID) (*service.BlogIndexSettings, error) {
			if inputAdminID != adminID {
				t.Fatalf("unexpected admin id: %s", inputAdminID)
			}
			received = input
			return &service.BlogIndexSettings{
				HeroBadge:       "CREATOR SPACE",
				HeroTitle:       "NanaMiku Blog",
				HeroDescription: "展示博客模板的首屏内容。",
				HeroActions: []service.BlogIndexHeroAction{
					{Label: "看最新文章", Href: "#latest-posts"},
				},
				QuickStats: []service.BlogIndexQuickStat{
					{Label: "模板栈", Value: "Astro + Vue + Go"},
				},
				FocusCard: service.BlogIndexFocusCard{
					Badge:       "本月在做什么",
					Title:       "打磨创作者模板",
					Description: "继续收口默认配置。",
					Footnote:    "文章来自后台管理面板发布",
				},
				ScrollCue: service.BlogIndexScrollCue{
					Label:     "向下阅读",
					AriaLabel: "向下阅读",
				},
			}, nil
		},
	}, logger)

	ctx := newAdminTestContext(
		consts.MethodPut,
		"/api/v1/admin/site-settings/blog-index",
		`{"hero_badge":" CREATOR SPACE ","hero_title":" NanaMiku Blog ","hero_description":" 展示博客模板的首屏内容。 ","hero_actions":[{"label":" 看最新文章 ","href":" #latest-posts "}],"quick_stats":[{"label":" 模板栈 ","value":" Astro + Vue + Go "}],"focus_card":{"badge":" 本月在做什么 ","title":" 打磨创作者模板 ","description":" 继续收口默认配置。 ","footnote":" 文章来自后台管理面板发布 "},"scroll_cue":{"label":" 向下阅读 ","aria_label":" 向下阅读 "}}`,
	)
	ctx.Set("admin_id", adminID)

	handler.UpdateBlogIndex(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusOK {
		t.Fatalf("unexpected status code: %d", ctx.Response.StatusCode())
	}

	if received.HeroTitle != " NanaMiku Blog " || len(received.HeroActions) != 1 {
		t.Fatalf("request payload was not bound as expected: %+v", received)
	}

	if len(logger.calls) != 1 || logger.calls[0].targetID != "blog_index" {
		t.Fatalf("unexpected audit log calls: %+v", logger.calls)
	}
}
