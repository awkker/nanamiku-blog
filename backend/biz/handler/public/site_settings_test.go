package public

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type stubSiteSettingsService struct {
	footer           *service.FooterSettings
	footerErr        error
	siteProfile      *service.SiteProfileSettings
	siteProfileErr   error
	homeHero         *service.HomeHeroSettings
	homeHeroErr      error
	homeAssets       *service.HomeAssetsSettings
	homeAssetsErr    error
	authorProfile    *service.AuthorProfileSettings
	authorProfileErr error
	integrations     *service.SiteIntegrationsSettings
	integrationsErr  error
}

func (s *stubSiteSettingsService) GetFooterSettings(context.Context) (*service.FooterSettings, error) {
	return s.footer, s.footerErr
}

func (s *stubSiteSettingsService) GetSiteProfileSettings(context.Context) (*service.SiteProfileSettings, error) {
	return s.siteProfile, s.siteProfileErr
}

func (s *stubSiteSettingsService) GetHomeHeroSettings(context.Context) (*service.HomeHeroSettings, error) {
	return s.homeHero, s.homeHeroErr
}

func (s *stubSiteSettingsService) GetHomeAssetsSettings(context.Context) (*service.HomeAssetsSettings, error) {
	return s.homeAssets, s.homeAssetsErr
}

func (s *stubSiteSettingsService) GetAuthorProfileSettings(context.Context) (*service.AuthorProfileSettings, error) {
	return s.authorProfile, s.authorProfileErr
}

func (s *stubSiteSettingsService) GetSiteIntegrationsSettings(context.Context) (*service.SiteIntegrationsSettings, error) {
	return s.integrations, s.integrationsErr
}

type responseEnvelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func newTestContext(method, uri string) *app.RequestContext {
	ctx := app.NewContext(0)
	ctx.Request.SetMethod(method)
	ctx.Request.SetRequestURI(uri)
	return ctx
}

func decodeResponse(t *testing.T, ctx *app.RequestContext) responseEnvelope {
	t.Helper()

	var body responseEnvelope
	if err := json.Unmarshal(ctx.Response.Body(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return body
}

func TestSiteSettingsHandlerGetHomeHeroSuccess(t *testing.T) {
	t.Parallel()

	handler := NewSiteSettingsHandler(&stubSiteSettingsService{
		homeHero: &service.HomeHeroSettings{
			HeroTitle:    "Miku Blog",
			HeroSubtitle: "Hello world",
		},
	})

	ctx := newTestContext(consts.MethodGet, "/api/v1/site-settings/home-hero")
	handler.GetHomeHero(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusOK {
		t.Fatalf("unexpected status code: %d", ctx.Response.StatusCode())
	}

	resp := decodeResponse(t, ctx)
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d", resp.Code)
	}

	var data service.HomeHeroSettings
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}

	if data.HeroTitle != "Miku Blog" || data.HeroSubtitle != "Hello world" {
		t.Fatalf("unexpected payload: %+v", data)
	}
}

func TestSiteSettingsHandlerGetSiteIntegrationsError(t *testing.T) {
	t.Parallel()

	handler := NewSiteSettingsHandler(&stubSiteSettingsService{
		integrationsErr: errors.New("boom"),
	})

	ctx := newTestContext(consts.MethodGet, "/api/v1/site-settings/site-integrations")
	handler.GetSiteIntegrations(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", ctx.Response.StatusCode())
	}

	resp := decodeResponse(t, ctx)
	if resp.Code == 0 {
		t.Fatalf("expected error response, got success")
	}
	if resp.Message != "failed to get site integrations settings" {
		t.Fatalf("unexpected response message: %q", resp.Message)
	}
}
