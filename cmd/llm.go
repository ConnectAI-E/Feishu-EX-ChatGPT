package main

import (
	"os"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/domain"
	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/repo/chatgpt"
	"github.com/sashabaranov/go-openai"
)

func newLLM() domain.LLMer {
	var (
		openAIToken = os.Getenv("OPENAI_TOKEN")

		// 如果不为空，则设置 openai 的代理模式
		openAIURL = os.Getenv("OPENAI_URL")
	)

	config := openai.DefaultConfig(openAIToken)
	if openAIURL != "" {
		config.BaseURL = openAIURL
	}

	openaiClient := openai.NewClientWithConfig(config)
	llm := chatgpt.NewChatGPT(openaiClient)

	return llm
}
