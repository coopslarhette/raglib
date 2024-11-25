package exa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the HTTP client for querying Exa API endpoints
type Client struct {
	apiKey string
	client *http.Client
}

type Contents struct {
	Text       *TextContent       `json:"text,omitempty"`
	Highlights *HighlightsContent `json:"highlights,omitempty"`
}

type TextContent struct {
	MaxCharacters   int  `json:"maxCharacters"`
	IncludeHtmlTags bool `json:"includeHtmlTags"`
}

type HighlightsContent struct {
	NumSentences     int    `json:"numSentences"`
	HighlightsPerUrl int    `json:"highlightsPerUrl"`
	Query            string `json:"query"`
}

type SummaryContent struct {
	Query string `json:"query"`
}

// SearchRequest represents the request structure for the Exa API search endpoint
type SearchRequest struct {
	Query              string     `json:"query"`
	UseAutoprompt      bool       `json:"useAutoprompt"`
	Type               string     `json:"type,omitempty"`
	Category           string     `json:"category,omitempty"`
	NumResults         int        `json:"numResults"`
	IncludeDomains     []string   `json:"includeDomains,omitempty"`
	ExcludeDomains     []string   `json:"excludeDomains,omitempty"`
	StartCrawlDate     *time.Time `json:"startCrawlDate,omitempty"`
	EndCrawlDate       *time.Time `json:"endCrawlDate,omitempty"`
	StartPublishedDate *time.Time `json:"startPublishedDate,omitempty"`
	EndPublishedDate   *time.Time `json:"endPublishedDate,omitempty"`
	IncludeText        []string   `json:"includeText,omitempty"`
	ExcludeText        []string   `json:"excludeText,omitempty"`
	Contents           *Contents  `json:"contents,omitempty"`
}

type SearchResponse struct {
	Results            []SearchResult `json:"results"`
	ResolvedSearchType string         `json:"resolvedSearchType,omitempty"`
	AutopromptString   string         `json:"autopromptString,omitempty"`
}

type SearchResult struct {
	Title           string    `json:"title,omitempty"`
	URL             string    `json:"url,omitempty"`
	PublishedDate   string    `json:"publishedDate,omitempty"`
	Author          string    `json:"author,omitempty"`
	Score           float64   `json:"score,omitempty"`
	ID              string    `json:"id,omitempty"`
	Text            string    `json:"text,omitempty"`
	Highlights      []string  `json:"highlights,omitempty"`
	HighlightScores []float64 `json:"highlightScores,omitempty"`
	Summary         string    `json:"summary,omitempty"`
}

func NewClient(apiKey string, client *http.Client) *Client {
	return &Client{
		apiKey: apiKey,
		client: client,
	}
}

func (c *Client) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.exa.ai/search", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error constructing request for Exa API: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while executing request to Exa API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is not OK; recieved code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Exa API response body: %v", err)
	}

	var result SearchResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing Exa API response: %v", err)
	}

	return &result, nil
}
