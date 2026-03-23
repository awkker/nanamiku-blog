package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"nanamiku-blog/backend/biz/bootstrap"
	"nanamiku-blog/backend/biz/service"
	"nanamiku-blog/backend/query"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	username := flag.String("username", "admin", "admin username")
	password := flag.String("password", "", "new password (or set ADMIN_NEW_PASSWORD)")
	force := flag.Bool("force", false, "skip password strength check")
	flag.Parse()

	uname := strings.TrimSpace(*username)
	if uname == "" {
		log.Fatal("username cannot be empty")
	}

	newPassword := strings.TrimSpace(*password)
	if newPassword == "" {
		newPassword = strings.TrimSpace(os.Getenv("ADMIN_NEW_PASSWORD"))
	}
	if newPassword == "" {
		log.Fatal("missing password: use -password or set ADMIN_NEW_PASSWORD")
	}

	if !*force {
		if err := validatePassword(newPassword); err != nil {
			log.Fatalf("weak password: %v (or pass -force)", err)
		}
	}

	cfg := bootstrap.LoadConfig()
	ctx := context.Background()

	db, err := bootstrap.NewDBPool(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	q := query.New(db)

	hash, err := service.HashPassword(newPassword)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	adminID, err := q.UpdateAdminPasswordByUsername(ctx, query.UpdateAdminPasswordByUsernameParams{
		Username:     uname,
		PasswordHash: hash,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Fatalf("admin user not found: %s", uname)
		}
		log.Fatalf("update password: %v", err)
	}

	if err := q.RevokeAllUserTokens(ctx, adminID); err != nil {
		log.Fatalf("password updated, but failed to revoke existing sessions: %v", err)
	}

	fmt.Printf("Password reset success for admin: %s\n", uname)
	fmt.Println("All existing refresh-token sessions have been revoked.")
}

func validatePassword(password string) error {
	runes := []rune(password)
	if len(runes) < 10 {
		return fmt.Errorf("must be at least 10 characters")
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
