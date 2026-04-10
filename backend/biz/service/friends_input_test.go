package service

import (
	"errors"
	"testing"
)

func TestNormalizeFriendLinkInput(t *testing.T) {
	result, err := normalizeFriendLinkInput(FriendLinkInput{
		Name:        "  Miku Blog Starter  ",
		Description: "  一个可继续 DIY 的亮色博客模板。  ",
		URL:         "example.com",
		AvatarURL:   "https://cdn.example.com/avatar.png",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "Miku Blog Starter" {
		t.Fatalf("unexpected name: %q", result.Name)
	}
	if result.Description != "一个可继续 DIY 的亮色博客模板。" {
		t.Fatalf("unexpected description: %q", result.Description)
	}
	if result.URL != "https://example.com" {
		t.Fatalf("unexpected url: %q", result.URL)
	}
	if result.Domain != "example.com" {
		t.Fatalf("unexpected domain: %q", result.Domain)
	}
	if result.AvatarURL != "https://cdn.example.com/avatar.png" {
		t.Fatalf("unexpected avatar url: %q", result.AvatarURL)
	}
}

func TestNormalizeFriendApplicationInput(t *testing.T) {
	result, err := normalizeFriendApplicationInput(FriendApplicationInput{
		SiteName:     "  Example Site  ",
		SiteURL:      "https://example.com/",
		AvatarURL:    "",
		Description:  "  A calm corner for notes. ",
		ContactEmail: " contact@example.com ",
		ContactNote:  "  已添加本站链接。  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.SiteName != "Example Site" {
		t.Fatalf("unexpected site name: %q", result.SiteName)
	}
	if result.SiteURL != "https://example.com" {
		t.Fatalf("unexpected site url: %q", result.SiteURL)
	}
	if result.ContactEmail != "contact@example.com" {
		t.Fatalf("unexpected email: %q", result.ContactEmail)
	}
	if result.ContactNote != "已添加本站链接。" {
		t.Fatalf("unexpected note: %q", result.ContactNote)
	}
}

func TestNormalizeFriendApplicationInputRejectsInvalidEmail(t *testing.T) {
	_, err := normalizeFriendApplicationInput(FriendApplicationInput{
		SiteName:     "Example Site",
		SiteURL:      "https://example.com",
		Description:  "Description",
		ContactEmail: "not-an-email",
	})
	if !errors.Is(err, ErrInvalidFriendInput) {
		t.Fatalf("expected ErrInvalidFriendInput, got %v", err)
	}
}
