package admin

import (
	"strings"
	"time"
)

func parseScheduledAtRFC3339(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, strings.TrimSpace(value))
}
