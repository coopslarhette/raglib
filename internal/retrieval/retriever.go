package retrieval

import (
	"context"
	"raglib/internal/document"
)

type Retriever interface {
	Query(ctx context.Context, query string, maxTopK uint64) ([]document.Document, error)
}
