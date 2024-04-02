package generation

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"raglib/internal/document"
	"strings"
)

var (
	promptTemplate = `Given the following document passages

<passages>%v</passages>

Use them to answer the text below. Your answer could take many forms depending on the text below. It could be as simple as a further elaboration on a basic understanding of a topic (using the info in the document), or it could be giving a detailed answer at a graduate level of explanation.

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

// Generate generates an answer to some text based off the given docs
// TODO Generators as a concept are feeling a light silly/like they haven't hit their mark ye
func (tg Answerer) Generate(ctx context.Context, seedInput string, documents []document.Document) (string, error) {
	combinedPassages := make([]string, len(documents))
	for i, d := range documents {
		combinedPassages[i] = documentToPassagesString(d)
	}
	passages := strings.Join(combinedPassages, "\n\n")

	prompt := fmt.Sprintf(promptTemplate, passages, seedInput)
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4TurboPreview,
		Temperature: 0,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}

	resp, err := tg.openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error making OpenAI API request: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func NewAnswerer(openAIClient *openai.Client) Answerer {
	return Answerer{openAIClient}
}
