package retrieval

import (
	"context"
	"github.com/coopslarhette/raglib/lib/document"
)

type Retriever interface {
	Query(ctx context.Context, query string, topK int) ([]document.Document, error)
}
