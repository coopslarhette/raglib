# raglib

raglib is a Go library for retrieval-augmented generation, providing a basic set of tools and abstractions for building applications that combine information retrieval and language model-based text generation. (Its currently a work in progress, thoughts and PRs welcome!)

## Features

- Retrieve relevant documents from various sources, such as web search results and vector databases
- Generate text based on the retrieved documents and user input
- Customize and extend the library components to fit specific use cases

### Examples

There's examples of it being used in a REST API contained in _./api_, that is then consumed by a NextJS + TypeScript web UI in _./web-client_ 

## Installation

To use raglib in your Go project, run:

```bash
go get github.com/coopslarhette/raglib
```

## Usage

### Retrieving Documents

raglib provides the `Retriever` interface for retrieving relevant documents based on a given query.

```go
import (
    "context"
    "github.com/coopslarhette/raglib/lib/document"
)

type Retriever interface {
    Query(ctx context.Context, query string, topK int) ([]document.Document, error)
}
```

The library includes two implementations of this interface:

1. `SERPRetriever`: Retrieves document snippets/web ranking by scraping Google Search results pages for a given query using the SERP API.
2. `ExaRetriever`: Retrieves full text of relevant web pages from https://exa.ai/, based on a given query (or URL)
3. `QdrantRetriever`: Retrieves relevant documents from a collections in a [Qdrant](https://qdrant.tech/) vector database, based on a given query.

An example of how to use the `SERPRetriever`:

```go
client := retrieval.NewSERPAPIClient("your_api_key")
retriever := retrieval.NewSERPRetriever(client)

query := "your search query"
topK := 10

docs, err := retriever.Query(context.Background(), query, topK)
if err != nil {
    // Handle error
}

for i, d := docs {
    fmt.Printf("document at position %d is title: %v", i, d.Title)
}
```

### Generating Text

raglib also provides the Generator interface for retrieving relevant documents based on a given query.:

```go
import (
	"context"
	"github.com/coopslarhette/raglib/lib/document"
)

type Generator interface {
    Generate(ctx context.Context, documents []document.Document, rawChunkChan chan<- string) error
}
```

There is one included implementation of the `Generator` interface: the `Answerer` struct for answering input text based on the retrieved documents and user input. The `Answerer` currently uses the OpenAI API for text generation. A facade that allows various model providers, or local LLMs to be used is forthcoming. 

Here's an example of how to use the `Answerer`:

```go
openaiClient := openai.NewClient("your_api_key")
answerer := generation.NewAnswerer(openaiClient)

seedInput := "user input for text generation"
documents := // Retrieved documents

rawChunkChan := make(chan string)
shouldStream := true

go answerer.Generate(ctx, prompt, documents, rawChunkChan, shouldStream)

// Consume the stream of generated text from the model provider (OpenAI in this case)
for response := range rawChunkChan {
    fmt.Print(response)
}
```

## Document Struct

The `document` package defines the `Document` struct, which represents a document retrieved by the `Retriever`. A `Document` consists of:

- `Passages`: A list of relevant passages from the document
- `Title`: The title of the document
- `Source`: The type of corpus the document came from (e.g., web, personal)
- `WebReference`: Information about the document's web source (if applicable)

## Contributing

Contributions to raglib are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request.
