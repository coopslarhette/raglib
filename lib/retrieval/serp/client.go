package serp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Client is the HTTP client for querying SerpApi API endpoints, mainly the
// Google Search endpoint.
type Client struct {
	apiKey string
	client *http.Client
}

type OrganicResult struct {
	Position         int    `json:"position"`
	Title            string `json:"title"`
	Link             string `json:"link"`
	Snippet          string `json:"snippet"`
	RedirectLink     string `json:"redirect_link"`
	DisplayedLink    string `json:"displayed_link"`
	Thumbnail        string `json:"thumbnail"`
	Date             string `json:"date"`
	Author           string `json:"author"`
	CitedBy          string `json:"cited_by"`
	ExtractedCitedBy int    `json:"extracted_cited_by"`
	Favicon          string `json:"favicon"`
}

// SearchResult represents the JSON response from SerpApi's API response. The fields
// here don't represent all the fields SerpApi returns, just the ones that might be
// interesting to us currently.
type SearchResult struct {
	OrganicResults []OrganicResult `json:"organic_results"`
}

func NewClient(apiKey string, client *http.Client) *Client {
	return &Client{
		apiKey: apiKey,
		client: client,
	}
}

func (c *Client) Query(ctx context.Context, query string, topK int) (*SearchResult, error) {
	apiURL, err := c.makeURL(query, topK)
	if err != nil {
		return nil, fmt.Errorf("error constructing URL for SERP API request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error constructing request for SERP API: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while executing request to SERP API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading SERP API response body: %v", err)
	}

	var result SearchResult
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing SERP API respone")
	}

	return &result, nil
}

func (c *Client) makeURL(query string, topK int) (*url.URL, error) {
	apiUrl, err := url.Parse("https://serpapi.com/search")
	if err != nil {
		return nil, fmt.Errorf("error parsing SERP API url: %v", err)
	}

	// Set the query parameters
	params := url.Values{}
	params.Set("q", query)
	params.Set("api_key", c.apiKey)
	params.Set("engine", "google")
	params.Set("num", strconv.Itoa(topK))
	apiUrl.RawQuery = params.Encode()

	return apiUrl, nil
}
