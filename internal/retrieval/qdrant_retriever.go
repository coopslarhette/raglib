package retrieval

import (
	"context"
	"fmt"
	qdrant "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"raglib/internal/document"
)

// QdrantRetriever implements the Retriever interface
type QdrantRetriever struct {
	pointsClient   qdrant.PointsClient
	openaiClient   *openai.Client
	collectionName string
}

func (qr QdrantRetriever) toQueryEmbedding(ctx context.Context, query string) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{query},
		Model: openai.AdaEmbeddingV2,
	}

	resp, err := qr.openaiClient.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating vector from query: %v", err)
	}

	return resp.Data[0].Embedding, nil
}

func (qr QdrantRetriever) Query(ctx context.Context, query string, maxTopK uint64) ([]document.Document, error) {
	qe, err := qr.toQueryEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error creating query emedding: %v", err)
	}

	unfilteredSearchResult, err := qr.pointsClient.Search(ctx, &qdrant.SearchPoints{
		CollectionName: qr.collectionName,
		Vector:         qe,
		Limit:          maxTopK,
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return nil, fmt.Errorf("error when searching points: %v", err)
	}

	docs := make([]document.Document, len(unfilteredSearchResult.Result))
	for i, r := range unfilteredSearchResult.Result {
		docs[i] = document.Document{
			Passages: []document.Passage{
				// TODO: maybe setup Query to accept a kind of parser as an argument to
				//   handle different search results types
				{Text: r.Payload["text"].GetStringValue()},
			},
		}
	}
	return docs, nil
}

func NewQdrantRetriever(pointsClient qdrant.PointsClient, openaiClient *openai.Client, collectionName string) QdrantRetriever {
	return QdrantRetriever{pointsClient, openaiClient, collectionName}
}
