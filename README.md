# raglib

raglib is a Go library for retrieval-augmented generation, providing a set of tools and abstractions for building applications that combine information retrieval and language model-based text generation.

## Features

- Retrieve relevant documents from various sources, such as web search results and vector databases
- Generate text based on the retrieved documents and user input
- Customize and extend the library components to fit specific use cases

## Installation

To use raglib in your Go project, run:

```bash
go get github.com/yourusername/raglib
```

## Usage

### Retrieving Documents

raglib provides the `Retriever` interface for retrieving relevant documents based on a given query.

```go
type Retriever interface {
	Query(ctx context.Context, query string, topK uint64) ([]document.Document, error)
}
```

The library includes two implementations of this interface:

1. `SERPRetriever`: Retrieves documents by scraping Google Search results pages for a given query using the SERP API.
2. `QdrantRetriever`: Retrieves documents from a Qdrant vector database using an `text-embedding-ada-002` for query embedding.

An example of how to use the `SERPRetriever`:

```go
client := retrieval.NewSERPAPIClient("your_api_key")
retriever := retrieval.NewSERPRetriever(client)

query := "your search query"
topK := uint64(10)

docs, err := retriever.Query(context.Background(), query, topK)
if err != nil {
    // Handle error
}

for i, d := docs {
	fmt.Printf("document at position %d is title: %v", i, d.Title)
}
```

### Generating Text

raglib provides the `Answerer` struct for generating text based on the retrieved documents and user input. The `Answerer` currently uses the OpenAI API for text generation. A facade that allows various model providers, or local LLMs to be used is forthcoming. 

Here's an example of how to use the `Answerer`:

```go
openaiClient := openai.NewClient("your_api_key")
answerer := generation.NewAnswerer(openaiClient)

seedInput := "user input for text generation"
documents := // Retrieved documents

responseChan := make(chan string)
shouldStream := true

go func() {
    if err := answerer.Generate(ctx, prompt, documents, responseChan, shouldStream); err != nil {
        // Handle error    
    }
}()

// Consume the generated text from the responseChan
for response := range responseChan {
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

Contributions to raglib are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on the [GitHub repository](https://github.com/yourusername/raglib).¬