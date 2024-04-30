package generation

import (
	"context"
	"raglib/internal/document"
)

type Generator interface {
	Generate(ctx context.Context, documents []document.Document, responseChan chan<- string) error
}
