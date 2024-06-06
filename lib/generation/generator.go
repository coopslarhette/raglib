package generation

import (
	"context"
	"raglib/lib/document"
)

type Generator interface {
	Generate(ctx context.Context, documents []document.Document, responseChan chan<- string) error
}
