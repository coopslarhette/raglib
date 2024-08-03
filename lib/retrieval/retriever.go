package retrieval

import (
	"context"
	"raglib/lib/document"
)

type Retriever interface {
	Query(ctx context.Context, query string, topK int) ([]document.Document, error)
}
