package exa

import (
	"context"
	"fmt"
	"raglib/lib/document"
)

// Retriever implements the retrieval.Retriever interface for the Exa search service. It retrieves web documents using their API endpoint.
type Retriever struct {
	client *Client
}

func (er Retriever) Query(ctx context.Context, query string, topK int) ([]document.Document, error) {
	request := SearchRequest{
		Query:      query,
		NumResults: topK,
		Contents: &Contents{
			Text: &TextContent{
				MaxCharacters: 1000,
			},
		},
		UseAutoprompt: true,
	}

	result, err := er.client.Search(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error querying Exa API: %v", err)
	}

	docs := make([]document.Document, len(result.Results))
	for i, r := range result.Results {
		docs[i] = document.Document{
			Passages: []document.Passage{
				{Text: r.Text},
			},
			Corpus: document.Web,
			WebReference: &document.WebReference{
				Title:         r.Title,
				Link:          r.URL,
				DisplayedLink: r.URL, // Exa API doesn't provide a separate displayed link
				Blurb:         r.Summary,
				Date:          r.PublishedDate,
				Author:        r.Author,
				APISource:     "exa",
			},
			Title: r.Title,
		}

		// Add highlights as additional passages
		// TODO: Don't do this for now, revisit
		//for _, highlight := range r.Highlights {
		//	docs[i].Passages = append(docs[i].Passages, document.Passage{Text: highlight})
		//}
	}

	return docs, nil
}

func NewRetriever(client *Client) Retriever {
	return Retriever{client}
}
