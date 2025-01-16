package urls

import (
	"fmt"
	"net/url"
	"strings"
)

// URLParts holds the parsed components of a URL
type URLParts struct {
	Domain    string
	TLD       string
	Subdomain string
}

// List of known multi-part TLDs
var multiPartTLDs = map[string]bool{
	"co.uk":  true,
	"co.in":  true,
	"com.au": true,
	"au.uk":  true,
	"co.nz":  true,
	"co.jp":  true,
	"co.kr":  true,
	"com.br": true,
	"com.cn": true,
}

// FullDomain returns the domain and TLD combined (e.g., "example.com")
func (u *URLParts) FullDomain() string {
	return u.Domain + "." + u.TLD
}

// isMultiPartTLD checks if the last parts of the domain form a known multi-part TLD
func isMultiPartTLD(parts []string) (string, bool) {
	if len(parts) < 2 {
		return "", false
	}

	// Check last two parts
	possibleTLD := strings.Join(parts[len(parts)-2:], ".")
	if multiPartTLDs[possibleTLD] {
		return possibleTLD, true
	}

	return "", false
}

// Parse extracts domain, TLD, and subdomain from a URL string
func Parse(urlStr string) (*URLParts, error) {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	parts := strings.Split(parsedURL.Hostname(), ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid hostname format")
	}

	result := &URLParts{}

	// Check for multi-part TLD
	if multiTLD, found := isMultiPartTLD(parts); found {
		// Remove the multi-part TLD parts from consideration
		parts = parts[:len(parts)-2]
		result.TLD = multiTLD

		switch len(parts) {
		case 0:
			return nil, fmt.Errorf("invalid hostname format")
		case 1:
			result.Domain = parts[0]
		default:
			// Always take the last remaining part as the domain
			result.Domain = parts[len(parts)-1]

			// Handle subdomains
			subdomainParts := parts[:len(parts)-1]
			if subdomainParts[0] == "www" {
				if len(subdomainParts) > 1 {
					result.Subdomain = strings.Join(subdomainParts[1:], ".")
				}
			} else {
				result.Subdomain = strings.Join(subdomainParts, ".")
			}
		}
	} else {
		// Handle regular TLDs
		result.TLD = parts[len(parts)-1]
		switch len(parts) {
		case 2:
			result.Domain = parts[0]
		case 3:
			if parts[0] == "www" {
				result.Domain = parts[1]
			} else {
				result.Subdomain = parts[0]
				result.Domain = parts[1]
			}
		default:
			result.Domain = parts[len(parts)-2]
			if parts[0] == "www" {
				result.Subdomain = strings.Join(parts[1:len(parts)-2], ".")
			} else {
				result.Subdomain = strings.Join(parts[:len(parts)-2], ".")
			}
		}
	}

	return result, nil
}
