package bootstrap

import (
	"strings"
	"testing"
	"time"
)

func TestValidateForServerDevelopmentAllowsDefaultSecrets(t *testing.T) {
	cfg := &Config{
		App: AppConfig{Env: "development"},
		JWT: JWTConfig{
			Secret:     "change-me-in-production",
			AccessTTL:  15 * time.Minute,
			RefreshTTL: 24 * time.Hour,
		},
		DB:    DBConfig{Password: "miku_secret"},
		Redis: RedisConfig{Password: "miku_redis"},
		Session: SessionConfig{
			CookieSecure:   false,
			CookieSameSite: "lax",
		},
	}

	if err := cfg.ValidateForServer(); err != nil {
		t.Fatalf("expected development config to pass, got %v", err)
	}
}

func TestValidateForServerProductionRejectsDangerousDefaults(t *testing.T) {
	cfg := &Config{
		App: AppConfig{Env: "production"},
		JWT: JWTConfig{
			Secret:     "change-me-in-production",
			AccessTTL:  15 * time.Minute,
			RefreshTTL: 24 * time.Hour,
		},
		DB:    DBConfig{Password: "miku_secret"},
		Redis: RedisConfig{Password: "miku_redis"},
		Session: SessionConfig{
			CookieSecure:   false,
			CookieSameSite: "lax",
		},
	}

	err := cfg.ValidateForServer()
	if err == nil {
		t.Fatal("expected production config to fail")
	}

	for _, want := range []string{
		"JWT_SECRET is using the default placeholder",
		"DB_PASSWORD is using the default development password",
		"REDIS_PASSWORD is using the default development password",
		"COOKIE_SECURE must be true in production",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected %q in error %q", want, err.Error())
		}
	}
}

func TestValidateForServerRejectsInvalidSessionAndTTL(t *testing.T) {
	cfg := &Config{
		App: AppConfig{Env: "development"},
		JWT: JWTConfig{
			Secret:     "dev-secret",
			AccessTTL:  time.Hour,
			RefreshTTL: 30 * time.Minute,
		},
		Session: SessionConfig{
			CookieSecure:   false,
			CookieSameSite: "sideways",
		},
	}

	err := cfg.ValidateForServer()
	if err == nil {
		t.Fatal("expected config validation to fail")
	}

	for _, want := range []string{
		"JWT_ACCESS_TTL must be shorter than JWT_REFRESH_TTL",
		"COOKIE_SAME_SITE must be one of lax, strict, none",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected %q in error %q", want, err.Error())
		}
	}
}
