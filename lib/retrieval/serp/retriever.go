package serp

import (
	"context"
	"fmt"
	"raglib/lib/document"
)

// Retriever implements the Retriever interface for the SERP API. SERP obtains documents and web ranking by scraping the relevant Google
// Search results page for a given query.
type Retriever struct {
	client *Client
}

func (sr Retriever) Query(ctx context.Context, query string, topK int) ([]document.Document, error) {
	result, err := sr.client.Query(ctx, query, topK)
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
			Corpus: document.Web,
			WebReference: &document.WebReference{
				Title:         r.Title,
				Link:          r.Link,
				DisplayedLink: r.DisplayedLink,
				Blurb:         r.Snippet,
				Date:          r.Date,
				Favicon:       r.Favicon,
				Author:        r.Author,
				Thumbnail:     r.Thumbnail,
				APISource:     "serp",
			},
			Title: r.Title,
		}
	}

	return docs, nil
}

func NewRetriever(client *Client) Retriever {
	return Retriever{client}
}
