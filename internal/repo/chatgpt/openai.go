package chatgpt

import (
	"context"
	"errors"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/domain"

	"github.com/sashabaranov/go-openai"
)

var _ domain.LLMer = &ChatGPT{}

type ChatGPT struct {
	model  string
	client *openai.Client
}

type ChatGPTOption func(c *ChatGPT)

func WithModel(model string) ChatGPTOption {
	return func(c *ChatGPT) {
		c.model = model
	}
}

func NewChatGPT(client *openai.Client, opts ...ChatGPTOption) *ChatGPT {

	chat := &ChatGPT{
		model:  openai.GPT3Dot5Turbo,
		client: client,
	}

	for _, opt := range opts {
		opt(chat)
	}

	return chat
}

func (c ChatGPT) Chat(ctx context.Context, messages []domain.LlmMessage) (*domain.LlmAnswer, error) {

	chatGPTMessages := c.makeChatGPTMessage(messages)

	return c.send(ctx, chatGPTMessages)
}

func (c ChatGPT) makeChatGPTMessage(messages []domain.LlmMessage) []openai.ChatCompletionMessage {

	chatGPTMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, m := range messages {
		chatGPTMessages = append(chatGPTMessages, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	return chatGPTMessages
}

func (c ChatGPT) send(ctx context.Context, messages []openai.ChatCompletionMessage) (*domain.LlmAnswer, error) {

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}

	if choices := resp.Choices; len(choices) == 0 {
		return nil, errors.New("got empty ChatGPT response")
	}

	answer := c.convertLlmAnswer(resp)
	return answer, nil
}

func (c ChatGPT) convertLlmAnswer(openaiResp openai.ChatCompletionResponse) *domain.LlmAnswer {

	choices := openaiResp.Choices[0]

	return &domain.LlmAnswer{
		Role:    choices.Message.Role,
		Content: choices.Message.Content,
	}
}
