package utils

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

var slugValid = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// SanitizeSlug normalizes user input to a URL-safe slug (lowercase, dashes, alnum)
// Returns sanitized slug and a boolean indicating validity.
func SanitizeSlug(in string) (string, bool) {
	s := strings.ToLower(strings.TrimSpace(in))
	s = strings.ReplaceAll(s, "_", "-")
	// collapse spaces to dashes
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) < 3 || len(s) > 40 {
		return s, false
	}
	if !slugValid.MatchString(s) {
		return s, false
	}
	return s, true
}

// GenerateShareSlug creates a short random slug easy to remember.
// Format: hex segment + '-' + hex segment (total ~ 9 chars).
func GenerateShareSlug() string {
	buf := make([]byte, 4)
	rand.Read(buf)
	first := hex.EncodeToString(buf[:2])
	second := hex.EncodeToString(buf[2:])
	return first + "-" + second
}
