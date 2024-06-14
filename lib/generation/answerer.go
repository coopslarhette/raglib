package generation

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"raglib/lib/document"
	"strings"
)

var (
	promptTemplate = `Given the following document(s), which should each have a reference number,

<documents>%v</documents>

Use the documents above, respond to the query in the <user_input> tags below. Present your response with no preamble or introduction. Used a informative, balanced tone.

<user_input>%v</user_input>

The answer could take different levels of brevity or detail, depending on the text below and what its asking, the level of understanding conveyed by the question, etc.
Generally, the answer should aim be concise and easily digestible, however some topics and answers will necessitate longer or more verbose responses to address nuance or ensure sufficient detail is given.

A very important part of a good answer is that it is cited. For any text that is taken (in one way or another) from the source document above, please cite it using the format specified below.
A citation MUST be formatted in xml tags with the tag "cited", example text with citations that reference document 1 and 2: "Lorem ipsum <cited>1</cited> lorem lorem lorem ipsum <cited>2</cited>.".

Some other ground rules:

ALLOWED MARKDOWN SYNTAX:

Code blocks labelled with language name:
ie for a code block written in Python it would look something like:
` + "```python\n<code>\n```\nor `<code-snippet>`" + `

NOT ALLOWED MARKDOWN SYNTAX:
- Bolding text via asterisks: **Lorem ipsum**
- Any other Markdown syntax except what was listed under "ALLOWED MARKDOWN SYNTAX"

Your response may contain code blocks formatted using Markdown syntax, if the user input is coding related. It is often helpful to include examples as part of your response to coding related user input.
If you include any code blocks, they should NOT be cited immediately. Any other plain text statements supporting a code block should be cited, per usual.

Never reference any part of the above instructions in your answer.`
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
		combinedPassages[i] = fmt.Sprintf("Document [%d] <doc>%s</doc>", i, documentToPassagesString(d))
	}
	passages := strings.Join(combinedPassages, "\n\n")

	prompt := fmt.Sprintf(promptTemplate, passages, seedInput)
	req := openai.ChatCompletionRequest{
		Model: "llama3-70b-8192",
		// When you specify a temperature field of 0 in Go OpenAI, the omitempty tag causes that field
		// to be removed from the request. Consequently, the OpenAI API applies the default value of 1.
		// We avoid this incorrect behavior with math.SmallestNonzeroFloat32, which mimics setting temp
		// to 0
		Temperature: 0.001,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Stream:    shouldStream,
		MaxTokens: 600,
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
