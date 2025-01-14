package generation

import (
	"context"
	"fmt"
	"raglib/lib/document"
	"raglib/lib/modelproviders"
	"strings"
)

var (
	promptTemplate = `You are an advanced AI assistant designed to provide accurate, well-cited responses to user queries based on given reference documents. Your task is to ingest and absorb the provided documents, and use them to answer the user's query in a clear, concise, and informative manner.

Here are the reference documents you should use to formulate your response. Note the documents are 0-indexed, please reference them in this way.:

<reference_documents>
%s
</reference_documents>

And here is the user's query:

<user_query>
%s
</user_query>

Instructions for formulating your response:

1. Response Formulation:
   Based on the documents, formulate a response that directly addresses the user's query. Your response should:
   - Be clear, concise, and informative.
   - Provide sufficient detail to address any nuances in the query.
   - Explicitly state any conflicts or complementary information.
   - Do not include any preamble or introduction.
   - Use an informative and balanced tone.

2. Citation Format:
   - For any information taken from the reference documents, it is important to reference the source material via citations.
   - Format: <cited>X</cited>, where X is the document number (originally in <reference_documents>) of the supporting source material.
   - Place citations directly after the statement they're supporting.
   - If citing multiple documents at a single citation location, use comma's to separate the numbers 
   - Example: "The study found a significant increase in productivity <cited>1</cited>. Also, major productivity gains were had from better tooling <cited>2,3</cited>."

3. Document References:
   - Do not directly mention or acknowledge the existence of the source documents in your answer.
   - Instead, use citations to support the statements in your response.

4. Response Length:
   - Aim for concise, easily digestible answers.
   - Provide more detail when necessary to ensure sufficient information.

5. Code Blocks:
   - Include code blocks when relevant to coding-related queries.
   - Do not cite code blocks directly, but cite any supporting statements related to the code block.
   - Label code blocks with the appropriate language name.
   - Example:
     ` + "```" + `python
def example_function():
pass
` + "```" +
		`6. HTML Character Entities
   - You may see some HTML Character Entities in the source documents
   - Do not copy these style of representation, instead rewrite them so a human can easily understand, ie "&lt;" should become "<"

Your entire response will be visible to the user. Therefore, present your thoughts and final response in a cohesive, professional manner.

Please begin.`
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
		combinedPassages[i] = fmt.Sprintf("Document [%d] <docucment>%s</docucment>", i, documentToPassagesString(d))
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
