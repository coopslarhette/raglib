package generation

import "raglib/internal/document"

type Generator interface {
	Generate(documents []document.Document) (string, error)
}
