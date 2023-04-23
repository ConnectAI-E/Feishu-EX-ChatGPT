package chatgpt

import (
	"context"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/domain"

	"github.com/sashabaranov/go-openai"
)

var _ domain.ChatGPTer = &ChatGPT{}

type ChatGPT struct {
	client *openai.Client
}

func NewChatGPT(client *openai.Client) *ChatGPT {

	return &ChatGPT{
		client: client,
	}
}

func (c ChatGPT) Chat(ctx context.Context, messages []domain.GPTMessage) error {

	return nil
}
