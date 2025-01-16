package generation

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"
)

func MakeJSONSafe(s string) string {
	// If empty, return as is
	if s == "" {
		return s
	}

	// First pass: clean obviously problematic characters
	cleaned := cleanControlChars(s)

	// Test if it's JSON safe
	if isJSONSafe(cleaned) {
		return cleaned
	}

	// If not safe, do aggressive cleaning
	return cleanAggressively(cleaned)
}

// cleanControlChars removes control characters and handles escaping
func cleanControlChars(s string) string {
	// Replace common problematic characters
	replacer := strings.NewReplacer(
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
		"\"", "\\\"",
		"\\", "\\\\",
	)
	s = replacer.Replace(s)

	// Remove control characters
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, s)
}

// cleanAggressively removes all non-printable ASCII characters
func cleanAggressively(s string) string {
	// Keep only printable ASCII characters
	reg := regexp.MustCompile("[^\x20-\x7E]+")
	return reg.ReplaceAllString(s, "")
}

// isJSONSafe tests if a string can be safely encoded as JSON
func isJSONSafe(s string) bool {
	// Try to marshal the string
	_, err := json.Marshal(s)
	return err == nil
}
