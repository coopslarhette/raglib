package retrieval

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// SERPAPIClient is the HTTP client for querying SerpApi API endpoints, mainly the
// Google Search endpoint.
type SERPAPIClient struct {
	apiKey string
	client *http.Client
}

// SearchResult represents the JSON response from SerpApi's API response. The fields
// here don't represent all the fields SerpApi returns, just the ones that might be
// interesting to us currently.
type SearchResult struct {
	OrganicResults []struct {
		Position int    `json:"position"`
		Title    string `json:"title"`
		Link     string `json:"link"`
		Snippet  string `json:"snippet"`
	} `json:"organic_results"`
}

func NewSerpApiClient(apiKey string, client *http.Client) *SERPAPIClient {
	return &SERPAPIClient{
		apiKey: apiKey,
		client: client,
	}
}

func (c *SERPAPIClient) Query(ctx context.Context, query string, topK uint64) (*SearchResult, error) {
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

func (c *SERPAPIClient) makeURL(query string, topK uint64) (*url.URL, error) {
	apiUrl, err := url.Parse("https://serpapi.com/search")
	if err != nil {
		return nil, fmt.Errorf("error parsing SERP API url: %v", err)
	}

	// Set the query parameters
	params := url.Values{}
	params.Set("q", query)
	params.Set("api_key", c.apiKey)
	params.Set("engine", "google")
	params.Set("num", strconv.FormatUint(topK, 10))
	apiUrl.RawQuery = params.Encode()

	return apiUrl, nil
}
