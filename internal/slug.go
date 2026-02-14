package internal

import (
	"regexp"
	"strings"
)

var (
	nonAlphanumeric  = regexp.MustCompile(`[^a-z0-9-]+`)
	multipleHyphens  = regexp.MustCompile(`-{2,}`)
)

// Slugify converts a title into a URL-friendly slug.
// It lowercases, replaces non-alphanumeric characters with hyphens,
// collapses consecutive hyphens, and trims leading/trailing hyphens.
func Slugify(title string) string {
	s := strings.ToLower(title)
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = multipleHyphens.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// NormalizeTags takes a comma-separated tag string and returns
// a slice of cleaned, lowercase, hyphen-separated tags.
func NormalizeTags(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	var tags []string
	for _, p := range parts {
		tag := Slugify(strings.TrimSpace(p))
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}
