package service

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

var (
	ErrInvalidFriendInput         = errors.New("invalid friend input")
	ErrFriendApplicationDuplicate = errors.New("friend application already exists")
	ErrFriendApplicationNotFound  = errors.New("friend application not found")
	ErrFriendApplicationProcessed = errors.New("friend application already reviewed")
)

type FriendLinkInput struct {
	Name        string
	Description string
	URL         string
	Domain      string
	AvatarURL   string
}

type FriendApplicationInput struct {
	SiteName     string
	SiteURL      string
	AvatarURL    string
	Description  string
	ContactEmail string
	ContactNote  string
}

func normalizeFriendLinkInput(input FriendLinkInput) (FriendLinkInput, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return FriendLinkInput{}, fmt.Errorf("%w: name required", ErrInvalidFriendInput)
	}

	normalizedURL, err := normalizeExternalURL(input.URL)
	if err != nil {
		return FriendLinkInput{}, fmt.Errorf("%w: invalid url", ErrInvalidFriendInput)
	}

	avatarURL := strings.TrimSpace(input.AvatarURL)
	if avatarURL != "" {
		avatarURL, err = normalizeExternalURL(avatarURL)
		if err != nil {
			return FriendLinkInput{}, fmt.Errorf("%w: invalid avatar_url", ErrInvalidFriendInput)
		}
	}

	domain := strings.TrimSpace(input.Domain)
	if domain == "" {
		domain = deriveDomainFromURL(normalizedURL)
	}

	return FriendLinkInput{
		Name:        name,
		Description: strings.TrimSpace(input.Description),
		URL:         normalizedURL,
		Domain:      domain,
		AvatarURL:   avatarURL,
	}, nil
}

func normalizeFriendApplicationInput(input FriendApplicationInput) (FriendApplicationInput, error) {
	linkInput, err := normalizeFriendLinkInput(FriendLinkInput{
		Name:        input.SiteName,
		Description: input.Description,
		URL:         input.SiteURL,
		AvatarURL:   input.AvatarURL,
	})
	if err != nil {
		return FriendApplicationInput{}, err
	}

	contactEmail := strings.TrimSpace(input.ContactEmail)
	if contactEmail == "" {
		return FriendApplicationInput{}, fmt.Errorf("%w: contact_email required", ErrInvalidFriendInput)
	}
	if _, err := mail.ParseAddress(contactEmail); err != nil {
		return FriendApplicationInput{}, fmt.Errorf("%w: invalid contact_email", ErrInvalidFriendInput)
	}

	return FriendApplicationInput{
		SiteName:     linkInput.Name,
		SiteURL:      linkInput.URL,
		AvatarURL:    linkInput.AvatarURL,
		Description:  linkInput.Description,
		ContactEmail: contactEmail,
		ContactNote:  strings.TrimSpace(input.ContactNote),
	}, nil
}

func deriveDomainFromURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
}

func normalizeExternalURL(raw string) (string, error) {
	candidate := strings.TrimSpace(raw)
	if candidate == "" {
		return "", fmt.Errorf("empty url")
	}
	if !strings.Contains(candidate, "://") {
		candidate = "https://" + candidate
	}

	parsed, err := url.Parse(candidate)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("unsupported scheme")
	}
	if parsed.Hostname() == "" {
		return "", fmt.Errorf("missing host")
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)
	if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
		parsed.Path = ""
	}

	return parsed.String(), nil
}
