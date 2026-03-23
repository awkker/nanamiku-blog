package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"nanamiku-blog/backend/query"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired or revoked")
	ErrTokenInvalid       = errors.New("invalid token")
	ErrAccountConflict    = errors.New("account conflict")
	ErrInvalidAccount     = errors.New("invalid account")
	ErrWeakPassword       = errors.New("weak password")
)

const DefaultAdminAvatarURL = "/picture/author.jpg"

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type AuthService struct {
	q   *query.Queries
	db  *pgxpool.Pool
	cfg JWTConfig
}

func NewAuthService(db *pgxpool.Pool, cfg JWTConfig) *AuthService {
	return &AuthService{
		q:   query.New(db),
		db:  db,
		cfg: cfg,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type AdminClaims struct {
	jwt.RegisteredClaims
	AdminID  uuid.UUID `json:"admin_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
}

type AdminPublicProfile struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type UpdateAccountInput struct {
	Username    string
	Email       string
	NewPassword string
}

func (s *AuthService) Login(ctx context.Context, identifier, password string) (*TokenPair, error) {
	admin, err := s.q.GetAdminByIdentifier(ctx, strings.TrimSpace(identifier))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("query admin: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	pair, err := s.generateTokenPair(ctx, admin.ID, admin.Username, admin.Role)
	if err != nil {
		return nil, err
	}

	_ = s.q.UpdateAdminLastLogin(ctx, admin.ID)
	return pair, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, rawRefresh string) (*TokenPair, error) {
	hash := hashToken(rawRefresh)

	rt, err := s.q.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("query refresh token: %w", err)
	}

	admin, err := s.q.GetAdminByID(ctx, rt.AdminUserID)
	if err != nil {
		return nil, fmt.Errorf("query admin: %w", err)
	}

	_ = s.q.RevokeRefreshToken(ctx, rt.ID)

	pair, err := s.generateTokenPair(ctx, admin.ID, admin.Username, admin.Role)
	if err != nil {
		return nil, err
	}

	return pair, nil
}

func (s *AuthService) Logout(ctx context.Context, rawRefresh string) error {
	if rawRefresh == "" {
		return nil
	}
	hash := hashToken(rawRefresh)
	rt, err := s.q.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return nil
	}
	return s.q.RevokeRefreshToken(ctx, rt.ID)
}

func (s *AuthService) LogoutAll(ctx context.Context, adminID uuid.UUID) error {
	return s.q.RevokeAllUserTokens(ctx, adminID)
}

func (s *AuthService) GetAdminInfo(ctx context.Context, adminID uuid.UUID) (*query.GetAdminByIDRow, error) {
	admin, err := s.q.GetAdminByID(ctx, adminID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	return &admin, nil
}

func (s *AuthService) GetPublicProfile(ctx context.Context) (*AdminPublicProfile, error) {
	profile, err := s.q.GetPrimaryAdminPublicProfile(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &AdminPublicProfile{
				Username:    "admin",
				DisplayName: "Admin",
				AvatarURL:   DefaultAdminAvatarURL,
			}, nil
		}
		return nil, err
	}

	displayName := strings.TrimSpace(profile.DisplayName)
	if displayName == "" {
		displayName = strings.TrimSpace(profile.Username)
	}
	if displayName == "" {
		displayName = "Admin"
	}

	avatarURL := strings.TrimSpace(profile.AvatarUrl)
	if avatarURL == "" {
		avatarURL = DefaultAdminAvatarURL
	}

	return &AdminPublicProfile{
		Username:    profile.Username,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, adminID uuid.UUID, displayName, avatarURL string) (*query.GetAdminByIDRow, error) {
	current, err := s.GetAdminInfo(ctx, adminID)
	if err != nil {
		return nil, err
	}

	nextDisplayName := strings.TrimSpace(displayName)
	if nextDisplayName == "" {
		nextDisplayName = current.Username
	}

	nextAvatarURL := strings.TrimSpace(avatarURL)
	if nextAvatarURL == "" {
		nextAvatarURL = DefaultAdminAvatarURL
	}

	if err := s.q.UpdateAdminProfile(ctx, query.UpdateAdminProfileParams{
		ID:          adminID,
		DisplayName: nextDisplayName,
		AvatarUrl:   nextAvatarURL,
	}); err != nil {
		return nil, fmt.Errorf("update admin profile: %w", err)
	}

	updated, err := s.GetAdminInfo(ctx, adminID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *AuthService) UpdateAccount(ctx context.Context, adminID uuid.UUID, input UpdateAccountInput) (*query.GetAdminByIDRow, bool, error) {
	current, err := s.GetAdminInfo(ctx, adminID)
	if err != nil {
		return nil, false, err
	}

	nextUsername := strings.TrimSpace(input.Username)
	if nextUsername == "" {
		return nil, false, fmt.Errorf("%w: username cannot be empty", ErrInvalidAccount)
	}

	nextEmail := strings.ToLower(strings.TrimSpace(input.Email))
	if nextEmail == "" || !strings.Contains(nextEmail, "@") {
		return nil, false, fmt.Errorf("%w: email format invalid", ErrInvalidAccount)
	}

	nextPassword := strings.TrimSpace(input.NewPassword)
	if nextPassword != "" {
		if err := validateAdminPassword(nextPassword); err != nil {
			return nil, false, fmt.Errorf("%w: %v", ErrWeakPassword, err)
		}
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, false, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	qtx := s.q.WithTx(tx)
	passwordChanged := false

	if current.Username != nextUsername || current.Email != nextEmail {
		if updateErr := qtx.UpdateAdminAccount(ctx, query.UpdateAdminAccountParams{
			ID:       adminID,
			Username: nextUsername,
			Email:    nextEmail,
		}); updateErr != nil {
			if isUniqueViolation(updateErr) {
				return nil, false, ErrAccountConflict
			}
			return nil, false, fmt.Errorf("update admin account: %w", updateErr)
		}
	}

	if nextPassword != "" {
		hash, hashErr := HashPassword(nextPassword)
		if hashErr != nil {
			return nil, false, fmt.Errorf("hash password: %w", hashErr)
		}
		if updateErr := qtx.UpdateAdminPasswordByID(ctx, query.UpdateAdminPasswordByIDParams{
			ID:           adminID,
			PasswordHash: hash,
		}); updateErr != nil {
			return nil, false, fmt.Errorf("update admin password: %w", updateErr)
		}
		passwordChanged = true
	}

	if passwordChanged || current.Username != nextUsername || current.Email != nextEmail {
		if revokeErr := qtx.RevokeAllUserTokens(ctx, adminID); revokeErr != nil {
			return nil, false, fmt.Errorf("revoke old sessions: %w", revokeErr)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, false, fmt.Errorf("commit tx: %w", err)
	}

	updated, err := s.GetAdminInfo(ctx, adminID)
	if err != nil {
		return nil, false, err
	}

	return updated, passwordChanged, nil
}

func (s *AuthService) ValidateAccessToken(tokenStr string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AdminClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func (s *AuthService) generateTokenPair(ctx context.Context, adminID uuid.UUID, username, role string) (*TokenPair, error) {
	now := time.Now()
	accessExp := now.Add(s.cfg.AccessTTL)

	claims := &AdminClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		AdminID:  adminID,
		Username: username,
		Role:     role,
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	rawRefresh := uuid.New().String()
	refreshHash := hashToken(rawRefresh)
	refreshExp := now.Add(s.cfg.RefreshTTL)

	_, err = s.q.CreateRefreshToken(ctx, query.CreateRefreshTokenParams{
		AdminUserID: adminID,
		TokenHash:   refreshHash,
		ExpiresAt:   refreshExp,
	})
	if err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresAt:    accessExp.Unix(),
	}, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func validateAdminPassword(password string) error {
	runes := []rune(password)
	if len(runes) < 6 {
		return fmt.Errorf("must be at least 6 characters")
	}

	hasLetter := false
	hasDigit := false
	for _, r := range runes {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return fmt.Errorf("must include both letters and numbers")
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "23505"
}
