package retrieval

import (
	"context"
	"fmt"
	"raglib/internal/document"
)

// SERPRetriever implements the Retriever interface. It retrieves documents by scraping the relevant Google
// Search results page for a given query.
type SERPRetriever struct {
	client *SERPAPIClient
}

func (sr SERPRetriever) Query(ctx context.Context, query string, maxTopK uint64) ([]document.Document, error) {
	result, err := sr.client.Query(ctx, query, maxTopK)
	if err != nil {
		return nil, fmt.Errorf("error querying SERP API: %v", err)
	}

	docs := make([]document.Document, len(result.OrganicResults))
	for i, r := range result.OrganicResults {
		docs[i] = document.Document{
			Passages: []document.Passage{
				// TODO: maybe setup Query to accept a kind of parser as an argument to
				//   handle different search results types
				{Text: r.Snippet},
			},
			Source: document.Web,
			WebReference: &document.WebReference{
				Title:         r.Title,
				Link:          r.Link,
				DisplayedLink: r.DisplayedLink,
				Snippet:       r.Snippet,
				Date:          r.Date,
				Favicon:       r.Favicon,
				Author:        r.Author,
				Thumbnail:     r.Thumbnail,
			},
			Title: r.Title,
		}
	}

	return docs, nil
}

func NewSERPRetriever(client *SERPAPIClient) SERPRetriever {
	return SERPRetriever{client}
}
