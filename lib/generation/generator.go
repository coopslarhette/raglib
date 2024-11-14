package generation

import (
	"context"
	"raglib/lib/document"
)

type Generator interface {
	// Generate should contain any business logic pertaining to how to shape the models response, ie prompting, tool use
	// etc
	Generate(ctx context.Context, documents []document.Document, responseChan chan<- string) error
}
