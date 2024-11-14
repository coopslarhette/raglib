package generation

import (
	"context"
	"fmt"
	"raglib/lib/document"
	"raglib/lib/modelproviders"
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
	modelProvider *modelproviders.Facade
}

// Generate implements the Generator interface. It generates an answer to some text grounded in the given documents.
func (tg Answerer) Generate(ctx context.Context, seedInput string, documents []document.Document, rawChunkChan chan<- string, shouldStream bool) error {
	defer close(rawChunkChan)
	combinedPassages := make([]string, len(documents))
	for i, d := range documents {
		combinedPassages[i] = fmt.Sprintf("Document [%d] <doc>%s</doc>", i, documentToPassagesString(d))
	}
	passages := strings.Join(combinedPassages, "\n\n")

	prompt := fmt.Sprintf(promptTemplate, passages, seedInput)
	req := modelproviders.GenerateRequest{
		Provider:     modelproviders.AnthropicProvider,
		Model:        "claude-3-5-haiku-20241022",
		Prompt:       prompt,
		ShouldStream: shouldStream,
		MaxTokens:    600,
	}

	return tg.modelProvider.Generate(ctx, req, rawChunkChan)
}

func NewAnswerer(modelProvider *modelproviders.Facade) Answerer {
	return Answerer{modelProvider}
}
