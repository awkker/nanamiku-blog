package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"nanamiku-blog/backend/biz/bootstrap"
	"nanamiku-blog/backend/biz/service"
	"nanamiku-blog/backend/query"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := bootstrap.LoadConfig()
	ctx := context.Background()

	db, err := bootstrap.NewDBPool(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	q := query.New(db)

	usernameFlag := flag.String("username", "", "admin username (or set ADMIN_SEED_USERNAME)")
	emailFlag := flag.String("email", "", "admin email (or set ADMIN_SEED_EMAIL)")
	passwordFlag := flag.String("password", "", "admin password (or set ADMIN_SEED_PASSWORD)")
	flag.Parse()

	username := strings.TrimSpace(firstNonEmpty(*usernameFlag, os.Getenv("ADMIN_SEED_USERNAME")))
	email := strings.TrimSpace(firstNonEmpty(*emailFlag, os.Getenv("ADMIN_SEED_EMAIL")))
	password := strings.TrimSpace(firstNonEmpty(*passwordFlag, os.Getenv("ADMIN_SEED_PASSWORD")))

	if username == "" || email == "" || password == "" {
		log.Fatal("seed requires explicit username, email, and password via flags or ADMIN_SEED_* env vars")
	}
	if err := service.ValidateAdminPassword(password); err != nil {
		log.Fatalf("invalid password: %v", err)
	}

	hash, err := service.HashPassword(password)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	row, err := q.CreateAdminUser(ctx, query.CreateAdminUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Role:         "admin",
	})
	if err != nil {
		log.Fatalf("create admin: %v", err)
	}

	fmt.Printf("Admin user created: id=%s username=%s email=%s\n", row.ID, username, email)
	fmt.Println("Initial password accepted and stored successfully.")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
