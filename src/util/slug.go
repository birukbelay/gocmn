package util

import (
	"github.com/oklog/ulid/v2"
	"regexp"
	"strings"
)

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
	//take the first 7 characters
	if len(slug) > 7 {
		slug = slug[:7]
	}

	//this will be sorted by name
	return slug + "_" + ulid.Make().String()
}
