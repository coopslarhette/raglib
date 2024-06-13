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

<passages>Document [0] <doc>The basic for loop has three components separated by semicolons: the init statement: executed before the first iteration; the condition expression: evaluated ...</doc>

Document [1] <doc>A range loop is used to iterate over the values in a slice, array or channel, the keys and values in a map, or the characters in a string.</doc>

Document [2] <doc>for is Go's only looping construct. Here are some basic types of for loops. package main. import "fmt". func main() {. The most basic type, with a single ...</doc>

Document [3] <doc>The range form of the for loop iterates over a slice or map. When ranging over a slice, two values are returned for each iteration. The first is the index ...</doc>

Document [4] <doc>Golang for loop. In Golang, we use the for loop to repeat a block of code until the specified condition is met. Here's the syntax of the for loop in Golang.</doc></passages>

Use them to respond to the text in the user_input tags below.

<user_input>concisely answer this question: <question>golang for loop</question></user_input>

The answer could take different levels of brevity or detail, depending on the text below and what its asking, the level of understanding conveyed, etc. It could be as simple as a further elaboration on a basic understanding of a topic (using the info in the document(s)), or it could be giving a detailed answer at a graduate level of explanation.

Generally, the answer should be concise and easily digestible, unless being more verbose is appropriate given the content in the document(s) and user inputted text. 

A very important part of a good answer is that it is cited. For any text that is taken (in one way or another) from the source document above, please cite it by referencing the provided document number. When you cite a reference, please do so by putting it in xml tags with the tag "cited", i.e. "Lorem ipsum <cited>1</cited> lorem lorem lorem ipsum <cited>2</cited>.".

Some ground rules:

ALLOWED MARKDOWN SYNTAX:

Code blocks:
` + "```<language>\n<code-to-be-rendered>\n```\nor `<code>`" + `

NOT ALLOWED MARKDOWN SYNTAX:
- Bolding text via asterisks: **Lorem ipsum**
- Any other Markdown syntax except what was listed under "ALLOWED MARKDOWN SYNTAX"

Answer in plain text. Your plain text may contain code blocks formatted using Markdown syntax, if the user input is coding related (ie it wouldn't be appropriate to include them as part of a general information query)'. 

If you include any code blocks, they should NOT be cited immediately. Any other plain text statements supporting a code block should be cited, per usual.

Each newline will be rendered individually, we are not rendering the entire answer as Markdown, so any double newlines could look strange.`
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
	fmt.Println(prompt)
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4Turbo,
		Temperature: 0,
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