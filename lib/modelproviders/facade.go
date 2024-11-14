package modelproviders

import (
	"context"
	"errors"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/sashabaranov/go-openai"
	"io"
	"strings"
)

// ModelProvider represents supported model providers
type ModelProvider string

const (
	OpenAIProvider    ModelProvider = "openai"
	AnthropicProvider ModelProvider = "anthropic"
	GroqProvider      ModelProvider = "groq"
)

// Facade handles interactions with different LLM providers
type Facade struct {
	// TODO: Exposing the openai client publicly like this is a hack to we can use their Embeddings endpoint ot maintain
	// the Qdrant retrieval functionality, fix after landing ModelProviderFacade
	OpenAIClient    *openai.Client
	anthropicClient *anthropic.Client
	// We access Groq's API thru OpenAI SDK, by changing some request params such as base URL
	groqClient *openai.Client
}

// NewFacade creates a new facade with configured clients
func NewFacade(openAIKey, anthropicKey string, groqKey string) *Facade {
	groqConfig := openai.DefaultConfig(groqKey)
	groqConfig.BaseURL = "https://api.groq.com/openai/v1"

	return &Facade{
		OpenAIClient:    openai.NewClient(openAIKey),
		anthropicClient: anthropic.NewClient(option.WithAPIKey(anthropicKey)),
		groqClient:      openai.NewClientWithConfig(groqConfig),
	}
}

// GenerateRequest contains the basic parameters needed for generation
type GenerateRequest struct {
	Provider     ModelProvider
	Model        string
	Prompt       string
	ShouldStream bool
	MaxTokens    int
}

// Generate handles completion requests for different LLM providers
// TODO rename
func (f *Facade) Generate(ctx context.Context, req GenerateRequest, rawChunkChan chan<- string) error {
	switch req.Provider {
	case OpenAIProvider:
		return f.handleOpenAI(ctx, req, rawChunkChan)
	case AnthropicProvider:
		return f.handleAnthropic(ctx, req, rawChunkChan)
	case GroqProvider:
		return f.handleGroq(ctx, req, rawChunkChan)
	default:
		return fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

func (f *Facade) handleOpenAI(ctx context.Context, req GenerateRequest, rawChunkChan chan<- string) error {
	openAIReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Temperature: 0.001,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: req.Prompt},
		},
		Stream:    req.ShouldStream,
		MaxTokens: req.MaxTokens,
	}

	if !req.ShouldStream {
		resp, err := f.OpenAIClient.CreateChatCompletion(ctx, openAIReq)
		if err != nil {
			return fmt.Errorf("error making OpenAI API request: %v", err)
		}
		rawChunkChan <- resp.Choices[0].Message.Content
		return nil
	}

	stream, err := f.OpenAIClient.CreateChatCompletionStream(ctx, openAIReq)
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
		rawChunkChan <- response.Choices[0].Delta.Content
	}

}

func (f *Facade) handleAnthropic(ctx context.Context, req GenerateRequest, rawChunkChan chan<- string) error {
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(req.Prompt)),
	}

	if !req.ShouldStream {
		message, err := f.anthropicClient.Messages.New(ctx, anthropic.MessageNewParams{
			Model:     anthropic.F(req.Model),
			MaxTokens: anthropic.F(int64(req.MaxTokens)),
			Messages:  anthropic.F(messages),
		})
		if err != nil {
			return fmt.Errorf("error making Anthropic API request: %v", err)
		}
		if len(message.Content) > 0 {
			rawChunkChan <- message.Content[0].Text
		}
		return nil
	}

	stream := f.anthropicClient.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(req.Model),
		MaxTokens: anthropic.F(int64(req.MaxTokens)),
		Messages:  anthropic.F(messages),
	})

	debugAcc := strings.Builder{}
	for stream.Next() {
		event := stream.Current()
		switch delta := event.Delta.(type) {
		case anthropic.ContentBlockDeltaEventDelta:
			if delta.Text != "" {
				rawChunkChan <- delta.Text
				debugAcc.WriteString(delta.Text)
			}
		}
	}

	fmt.Println("hello")
	fmt.Println(debugAcc.String())

	if err := stream.Err(); err != nil {
		return fmt.Errorf("error while streaming Anthropic response: %v", err)
	}

	return nil
}

// We access Groq's API thru OpenAI SDK, by changing some request params such as base URL
func (f *Facade) handleGroq(ctx context.Context, req GenerateRequest, rawChunkChan chan<- string) error {
	groqReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Temperature: 0.001,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: req.Prompt},
		},
		Stream:    req.ShouldStream,
		MaxTokens: req.MaxTokens,
	}

	if !req.ShouldStream {
		resp, err := f.groqClient.CreateChatCompletion(ctx, groqReq)
		if err != nil {
			return fmt.Errorf("error making Groq API request: %v", err)
		}
		rawChunkChan <- resp.Choices[0].Message.Content
		return nil
	}

	stream, err := f.groqClient.CreateChatCompletionStream(ctx, groqReq)
	if err != nil {
		return fmt.Errorf("error making Groq API request: %v", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error while streaming Groq response: %v", err)
		}
		rawChunkChan <- response.Choices[0].Delta.Content
	}
}
