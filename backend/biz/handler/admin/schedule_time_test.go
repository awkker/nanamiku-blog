package admin

import (
	"testing"
	"time"
)

func TestParseScheduledAtRFC3339(t *testing.T) {
	parsed, err := parseScheduledAtRFC3339(" 2026-04-03T12:34:56Z ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2026, 4, 3, 12, 34, 56, 0, time.UTC)
	if !parsed.Equal(expected) {
		t.Fatalf("unexpected parsed time: %s", parsed)
	}
}

func TestParseScheduledAtRFC3339RejectsInvalidValue(t *testing.T) {
	if _, err := parseScheduledAtRFC3339("tomorrow"); err == nil {
		t.Fatal("expected invalid RFC3339 value to fail")
	}
}

func TestParseMomentPublishDefaultsToPublished(t *testing.T) {
	status, scheduledAt, ok := parseMomentPublish("", "")
	if !ok {
		t.Fatal("expected default publish status to succeed")
	}
	if status != "published" {
		t.Fatalf("expected published, got %q", status)
	}
	if scheduledAt != nil {
		t.Fatalf("expected no scheduled time, got %v", scheduledAt)
	}
}

func TestParseMomentPublishScheduledRequiresRFC3339(t *testing.T) {
	status, scheduledAt, ok := parseMomentPublish("scheduled", "bad-input")
	if ok {
		t.Fatalf("expected parse to fail, got status=%q time=%v", status, scheduledAt)
	}
}

func TestParseMomentPublishScheduledParsesTimestamp(t *testing.T) {
	status, scheduledAt, ok := parseMomentPublish("scheduled", "2026-04-03T12:34:56Z")
	if !ok {
		t.Fatal("expected scheduled parse to succeed")
	}
	if status != "scheduled" {
		t.Fatalf("expected scheduled, got %q", status)
	}
	if scheduledAt == nil || !scheduledAt.Equal(time.Date(2026, 4, 3, 12, 34, 56, 0, time.UTC)) {
		t.Fatalf("unexpected scheduled time: %v", scheduledAt)
	}
}
