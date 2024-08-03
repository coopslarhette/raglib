package retrieval

import (
	"context"
	"fmt"
	"raglib/lib/document"
	"raglib/lib/retrieval/exa"
)

// ExaRetriever implements the Retriever interface. It retrieves documents using the Exa Search API endpoint.
type ExaRetriever struct {
	client *exa.Client
}

func (er ExaRetriever) Query(ctx context.Context, query string, topK int) ([]document.Document, error) {
	request := exa.SearchRequest{
		Query:      query,
		NumResults: topK,
		Contents: &exa.Contents{
			Text: &exa.TextContent{
				MaxCharacters: 1000,
			},
		},
		UseAutoprompt: true,
	}

	fmt.Printf("%+v", request)

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
			Source: document.Web,
			WebReference: &document.WebReference{
				Title:         r.Title,
				Link:          r.URL,
				DisplayedLink: r.URL, // Exa API doesn't provide a separate displayed link
				Blurb:         r.Summary,
				Date:          r.PublishedDate,
				Author:        r.Author,
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

func NewExaRetriever(client *exa.Client) ExaRetriever {
	return ExaRetriever{client}
}
