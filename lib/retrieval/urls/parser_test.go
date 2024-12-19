package urls

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		wantDomain    string
		wantTLD       string
		wantSubdomain string
		wantErr       bool
	}{
		// Basic cases
		{
			name:       "simple domain",
			url:        "example.com",
			wantDomain: "example",
			wantTLD:    "com",
		},
		{
			name:       "www subdomain",
			url:        "www.example.com",
			wantDomain: "example",
			wantTLD:    "com",
		},
		{
			name:          "custom subdomain",
			url:           "blog.example.com",
			wantDomain:    "example",
			wantTLD:       "com",
			wantSubdomain: "blog",
		},

		// Multi-part TLD cases
		{
			name:          "complex multi-part TLD with www",
			url:           "www.mail.yahoo.co.in",
			wantDomain:    "yahoo",
			wantTLD:       "co.in",
			wantSubdomain: "mail",
		},
		{
			name:       "UK multi-part domain",
			url:        "www.abc.au.uk",
			wantDomain: "abc",
			wantTLD:    "au.uk",
		},
		{
			name:       "co.uk domain",
			url:        "http://www.google.co.uk",
			wantDomain: "google",
			wantTLD:    "co.uk",
		},

		// URL scheme cases
		{
			name:       "https scheme",
			url:        "https://github.com",
			wantDomain: "github",
			wantTLD:    "com",
		},
		{
			name:       "http scheme with country TLD",
			url:        "http://github.ca",
			wantDomain: "github",
			wantTLD:    "ca",
		},
		{
			name:       "https with www and country TLD",
			url:        "https://www.google.ru",
			wantDomain: "google",
			wantTLD:    "ru",
		},

		// Various formats
		{
			name:       "www prefix without scheme",
			url:        "www.yandex.com",
			wantDomain: "yandex",
			wantTLD:    "com",
		},
		{
			name:       "simple country TLD",
			url:        "yandex.ru",
			wantDomain: "yandex",
			wantTLD:    "ru",
		},

		// Error cases
		{
			name:    "single word domain",
			url:     "yandex",
			wantErr: true,
		},
		{
			name:    "empty string",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid URL format",
			url:     "http://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Domain != tt.wantDomain {
				t.Errorf("Domain = %v, want %v", got.Domain, tt.wantDomain)
			}
			if got.TLD != tt.wantTLD {
				t.Errorf("TLD = %v, want %v", got.TLD, tt.wantTLD)
			}
			if got.Subdomain != tt.wantSubdomain {
				t.Errorf("Subdomain = %v, want %v", got.Subdomain, tt.wantSubdomain)
			}
		})
	}
}

func TestURLParts_FullDomain(t *testing.T) {
	u := &URLParts{
		Domain: "example",
		TLD:    "com",
	}
	if got := u.FullDomain(); got != "example.com" {
		t.Errorf("FullDomain() = %v, want example.com", got)
	}
}
