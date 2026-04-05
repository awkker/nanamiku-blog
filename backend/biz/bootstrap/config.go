package bootstrap

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App     AppConfig
	Server  ServerConfig
	DB      DBConfig
	Redis   RedisConfig
	JWT     JWTConfig
	Session SessionConfig
	CORS    CORSConfig
	Weather WeatherConfig
	GeoIP   GeoIPConfig
	GitHub  GitHubConfig
}

type AppConfig struct {
	Env string
}

func (c AppConfig) IsProduction() bool {
	return strings.EqualFold(strings.TrimSpace(c.Env), "production")
}

type WeatherConfig struct {
	Location string
}

type GeoIPConfig struct {
	DBPath string
}

type GitHubConfig struct {
	APIToken   string
	APIBaseURL string
	CacheTTL   time.Duration
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c DBConfig) DSN() string {
	return "postgres://" + c.User + ":" + c.Password +
		"@" + c.Host + ":" + c.Port + "/" + c.Name +
		"?sslmode=" + c.SSLMode
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func (c RedisConfig) Addr() string {
	return c.Host + ":" + c.Port
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type SessionConfig struct {
	CookieDomain   string
	CookieSecure   bool
	CookieSameSite string
}

type CORSConfig struct {
	Origins []string
}

func LoadConfig() *Config {
	app := AppConfig{
		Env: envStr("APP_ENV", "development"),
	}

	return &Config{
		App: app,
		Server: ServerConfig{
			Port: envStr("SERVER_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     envStr("DB_HOST", "localhost"),
			Port:     envStr("DB_PORT", "5432"),
			User:     envStr("DB_USER", "miku"),
			Password: envStr("DB_PASSWORD", "miku_secret"),
			Name:     envStr("DB_NAME", "miku_blog"),
			SSLMode:  envStr("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     envStr("REDIS_HOST", "localhost"),
			Port:     envStr("REDIS_PORT", "6379"),
			Password: envStr("REDIS_PASSWORD", "miku_redis"),
			DB:       envInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     envStr("JWT_SECRET", "change-me-in-production"),
			AccessTTL:  envDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL: envDuration("JWT_REFRESH_TTL", 30*24*time.Hour),
		},
		Session: SessionConfig{
			CookieDomain:   envStr("COOKIE_DOMAIN", ""),
			CookieSecure:   envBool("COOKIE_SECURE", app.IsProduction()),
			CookieSameSite: envStr("COOKIE_SAME_SITE", "lax"),
		},
		CORS: CORSConfig{
			Origins: strings.Split(envStr("CORS_ORIGINS", "http://localhost:4321"), ","),
		},
		Weather: WeatherConfig{
			Location: envStr("WEATHER_LOCATION", "Shanghai"),
		},
		GeoIP: GeoIPConfig{
			DBPath: envStr("GEOIP_DB_PATH", ""),
		},
		GitHub: GitHubConfig{
			APIToken:   envStr("GITHUB_TOKEN", ""),
			APIBaseURL: envStr("GITHUB_API_BASE_URL", "https://api.github.com"),
			CacheTTL:   envDuration("GITHUB_CACHE_TTL", 30*time.Minute),
		},
	}
}

func (c *Config) ValidateForServer() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	var issues []string
	if c.JWT.AccessTTL <= 0 {
		issues = append(issues, "JWT_ACCESS_TTL must be greater than 0")
	}
	if c.JWT.RefreshTTL <= 0 {
		issues = append(issues, "JWT_REFRESH_TTL must be greater than 0")
	}
	if c.JWT.AccessTTL >= c.JWT.RefreshTTL {
		issues = append(issues, "JWT_ACCESS_TTL must be shorter than JWT_REFRESH_TTL")
	}

	sameSite := strings.ToLower(strings.TrimSpace(c.Session.CookieSameSite))
	switch sameSite {
	case "lax", "strict", "none":
	default:
		issues = append(issues, "COOKIE_SAME_SITE must be one of lax, strict, none")
	}

	if !c.App.IsProduction() {
		if len(issues) > 0 {
			return fmt.Errorf("%s", strings.Join(issues, "; "))
		}
		return nil
	}

	if c.JWT.Secret == "change-me-in-production" {
		issues = append(issues, "JWT_SECRET is using the default placeholder")
	}
	if c.DB.Password == "miku_secret" {
		issues = append(issues, "DB_PASSWORD is using the default development password")
	}
	if c.Redis.Password == "miku_redis" {
		issues = append(issues, "REDIS_PASSWORD is using the default development password")
	}
	if !c.Session.CookieSecure {
		issues = append(issues, "COOKIE_SECURE must be true in production")
	}

	if len(issues) > 0 {
		return fmt.Errorf("%s", strings.Join(issues, "; "))
	}
	return nil
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
