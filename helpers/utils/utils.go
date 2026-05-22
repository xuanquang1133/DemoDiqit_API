package utils

import (
	"regexp"
	"strings"
)

func Slug(raw string) string {
	if raw == "" {
		return ""
	}
	cleaned := strings.ToLower(strings.TrimSpace(raw))
	cleaned = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(cleaned, "-")
	cleaned = regexp.MustCompile(`-+`).ReplaceAllString(cleaned, "-")
	return strings.Trim(cleaned, "-")
}

func SKU(raw string) string {
	if raw == "" {
		return ""
	}
	cleaned := strings.ToUpper(strings.TrimSpace(raw))
	cleaned = regexp.MustCompile(`[^A-Z0-9]+`).ReplaceAllString(cleaned, "-")
	cleaned = regexp.MustCompile(`-+`).ReplaceAllString(cleaned, "-")
	return strings.Trim(cleaned, "-")
}
