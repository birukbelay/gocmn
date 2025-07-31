package util

import (
	"regexp"
	"strings"

	"github.com/oklog/ulid/v2"
)

// CreateSlug lowercases, trims space to -, removes special characters, returns first 7 char +_+ ulid
func CreateSlug(input string) string {
	// Convert to lowercase
	slug := strings.ToLower(input)
	// Trim spaces from the beginning and end
	slug = strings.TrimSpace(slug)
	// Replace consecutive spaces with a single hyphen
	slug = regexp.MustCompile("\\s+").ReplaceAllString(slug, "_")
	// Replace spaces with hyphens
	//slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters using a regular expression
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")
	//take the first 7 characters TODO: get options, to decide length
	if len(slug) > 7 {
		slug = slug[:7]
	}
	if slug == "" {
		return ulid.Make().String()
	}

	//this will be sorted by name
	return slug + "_" + ulid.Make().String()
}
