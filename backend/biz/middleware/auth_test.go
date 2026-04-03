package middleware

import "testing"

func TestResolveAccessTokenPrefersBearerHeader(t *testing.T) {
	token := resolveAccessToken("Bearer header-token", "cookie-token")
	if token != "header-token" {
		t.Fatalf("expected bearer token, got %q", token)
	}
}

func TestResolveAccessTokenFallsBackToCookie(t *testing.T) {
	token := resolveAccessToken("", "cookie-token")
	if token != "cookie-token" {
		t.Fatalf("expected cookie token, got %q", token)
	}
}

func TestResolveAccessTokenIgnoresBlankBearer(t *testing.T) {
	token := resolveAccessToken("Bearer   ", "cookie-token")
	if token != "cookie-token" {
		t.Fatalf("expected cookie fallback, got %q", token)
	}
}
