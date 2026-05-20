package normalize

import (
	"regexp"
	"strings"
)

// Slug normalizes a slug value from user input.
// Rules: lowercase, trim, replace any run of non-alphanumeric chars with a single dash,
// collapse consecutive dashes, remove leading/trailing dashes.
func Slug(raw string) string {
	if raw == "" {
		return ""
	}
	cleaned := strings.ToLower(strings.TrimSpace(raw))

	// Replace any run of non-alphanumeric characters with a single dash
	cleaned = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(cleaned, "-")

	// Collapse multiple consecutive dashes into one
	cleaned = regexp.MustCompile(`-+`).ReplaceAllString(cleaned, "-")

	// Remove leading and trailing dashes
	return strings.Trim(cleaned, "-")
}

// SKU normalizes a SKU value from user input.
// Rules: uppercase, trim, replace any run of non-alphanumeric chars with a single dash,
// collapse consecutive dashes, remove leading/trailing dashes.
func SKU(raw string) string {
	if raw == "" {
		return ""
	}
	cleaned := strings.ToUpper(strings.TrimSpace(raw))

	// Replace any run of non-alphanumeric characters with a single dash
	cleaned = regexp.MustCompile(`[^A-Z0-9]+`).ReplaceAllString(cleaned, "-")

	// Collapse multiple consecutive dashes into one
	cleaned = regexp.MustCompile(`-+`).ReplaceAllString(cleaned, "-")

	// Remove leading and trailing dashes
	return strings.Trim(cleaned, "-")
}
