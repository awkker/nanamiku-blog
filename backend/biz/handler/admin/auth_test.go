package admin

import "testing"

func TestResolveRefreshTokenPrefersRequestBody(t *testing.T) {
	token := resolveRefreshToken("request-token", "cookie-token")
	if token != "request-token" {
		t.Fatalf("expected request token, got %q", token)
	}
}

func TestResolveRefreshTokenFallsBackToCookie(t *testing.T) {
	token := resolveRefreshToken("   ", "cookie-token")
	if token != "cookie-token" {
		t.Fatalf("expected cookie token, got %q", token)
	}
}

func TestResolveRefreshTokenReturnsEmptyWhenMissing(t *testing.T) {
	token := resolveRefreshToken("", "")
	if token != "" {
		t.Fatalf("expected empty token, got %q", token)
	}
}
