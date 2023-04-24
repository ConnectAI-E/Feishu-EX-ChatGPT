package domain

import "context"

type LlmMessage struct {
	Role    string
	Content string
}

type LlmAnswer struct {
	Role    string
	Content string
}

type LLMer interface {
	Chat(ctx context.Context, messages []LlmMessage) (*LlmAnswer, error)
}
