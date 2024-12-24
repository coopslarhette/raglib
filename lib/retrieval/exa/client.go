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

const (
	exaAPIBaseURL = "https://api.exa.com"
)

func buildExaAPIURL(endpoint string) string {
	return exaAPIBaseURL + endpoint
}

// Client is the HTTP client for querying Exa API endpoints
type Client struct {
	apiKey string
	client *http.Client
}

func NewClient(apiKey string, client *http.Client) *Client {
	return &Client{
		apiKey: apiKey,
		client: client,
	}
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

func (c *Client) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, buildExaAPIURL("/search"), bytes.NewBuffer(requestBody))
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

// ContentsRequest represents the request structure for the Exa API contents endpoint
type ContentsRequest struct {
	IDs  []string     `json:"ids"`
	Text *TextContent `json:"text,omitempty"`
	// TODO: consider renaming to HighlightsSpecification or something similar
	Highlights *HighlightsContent `json:"highlights,omitempty"`
	// TODO: same as above
	Summary *SummaryContent `json:"summary,omitempty"`
}

// ContentsResponse represents the response structure from the Exa API contents endpoint
type ContentsResponse struct {
	Results []ContentResult `json:"results"`
}

// ContentResult represents a single result in the contents response
type ContentResult struct {
	ID         string   `json:"id"`
	Text       string   `json:"text,omitempty"`
	Highlights []string `json:"highlights,omitempty"`
	Summary    string   `json:"summary,omitempty"`
}

// Contents retrieves content details for specified IDs from the Exa API
func (c *Client) Contents(ctx context.Context, request ContentsRequest) (*ContentsResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, buildExaAPIURL("/contents"), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error constructing URL contents request for Exa API: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while executing request to Exa API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is not OK; received code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Exa API response body: %v", err)
	}

	var result ContentsResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing Exa API response: %v", err)
	}

	return &result, nil
}
