package main

import (
	"context"
	"os"

	"github.com/agi-cn/llmplugin"
	"github.com/agi-cn/llmplugin/llm/openai"
	"github.com/agi-cn/llmplugin/plugins/calculator"
	"github.com/agi-cn/llmplugin/plugins/google"
)

func newLLMPluginManager() *llmplugin.PluginManager {
	var (
		openAIToken = os.Getenv("OPENAI_TOKEN")
	)

	chatgpt := openai.NewChatGPT(openAIToken)

	plugins := makePlugins()

	return llmplugin.NewPluginManager(
		chatgpt,
		llmplugin.WithPlugins(plugins),
	)
}

func makePlugins() []llmplugin.Plugin {
	var (
		googleEngineID = os.Getenv("GOOGLE_ENGINE_ID")
		googleToken    = os.Getenv("GOOGLE_TOKEN")
	)

	plugins := []llmplugin.Plugin{
		&llmplugin.SimplePlugin{
			Name:         "Weather",
			InputExample: ``,
			Desc:         "Can check the weather forecast",
			DoFunc: func(ctx context.Context, query string) (answer string, err error) {
				answer = "Here is dummy weather plugin"
				return
			},
		},

		calculator.NewCalculator(),
	}

	if googleEngineID != "" && googleToken != "" {
		plugins = append(plugins,
			google.NewGoogle(googleEngineID, googleToken))
	}

	return plugins
}
