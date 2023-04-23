package domain

import "context"

type GPTMessage struct {
	Role    string
	Content string
}

type ChatGPTer interface {
	Chat(ctx context.Context, messages []GPTMessage) error
}
