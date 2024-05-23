package generation

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"raglib/internal/document"
	"strings"
)

var (
	promptTemplate = `Given the following document passages, which should each have a reference number,

<passages>%v</passages>

Use them to answer the text below. Your answer could take many forms depending on the text below. It could be as simple as a further elaboration on a basic understanding of a topic (using the info in the document), or it could be giving a detailed answer at a graduate level of explanation.

A very important part of a good answer is that it is cited. For each statement, or component of your answer, please cite it by referencing the provided reference number. When you cite a reference, please do so by putting it in xml tags with the tag "cited", i.e. "Lorem ipsum <cited>1</cited> lorem lorem lorem ipsum <cited>2</cited>.".

<text_to_answer>%v</text_to_answer>
`
)

func documentToPassagesString(doc document.Document) string {
	texts := make([]string, len(doc.Passages))
	for i, passage := range doc.Passages {
		texts[i] = passage.Text
	}
	return strings.Join(texts, "")
}

// Answerer implements the Generator interface.
type Answerer struct {
	// TODO abstraction here instead to allow different model providers
	openaiClient *openai.Client
}

// Generate implements the Generator interface. It generates an answer to some text grounded in the given documents.
// TODO Generators as a concept are feeling a light silly/like they haven't hit their mark ye
func (tg Answerer) Generate(ctx context.Context, seedInput string, documents []document.Document, responseChan chan<- string, shouldStream bool) error {
	defer close(responseChan)
	combinedPassages := make([]string, len(documents))
	for i, d := range documents {
		combinedPassages[i] = fmt.Sprintf("Document Reference Number [%d] %s", i, documentToPassagesString(d))
	}
	passages := strings.Join(combinedPassages, "\n\n")

	prompt := fmt.Sprintf(promptTemplate, passages, seedInput)
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		Temperature: 0,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Stream: shouldStream,
	}

	if !shouldStream {
		resp, err := tg.openaiClient.CreateChatCompletion(ctx, req)
		if err != nil {
			return fmt.Errorf("error making OpenAI API request: %v", err)
		}
		responseChan <- resp.Choices[0].Message.Content
		return nil
	}

	stream, err := tg.openaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("error making OpenAI API request: %v", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("error while streaming response: %v", err)
		}

		responseChan <- response.Choices[0].Delta.Content
	}
}

func NewAnswerer(openAIClient *openai.Client) Answerer {
	return Answerer{openAIClient}
}
