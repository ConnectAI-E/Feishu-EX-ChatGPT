package main

import (
	"os"

	"github.com/agi-cn/llmplugin"
	"github.com/agi-cn/llmplugin/llm/openai"
	"github.com/agi-cn/llmplugin/plugins/agicn_search"
	"github.com/agi-cn/llmplugin/plugins/calculator"
	"github.com/agi-cn/llmplugin/plugins/google"
	"github.com/agi-cn/llmplugin/plugins/stablediffusion"
	"github.com/sirupsen/logrus"
)

func newLLMPluginManager() *llmplugin.PluginManager {
	var (
		openAIToken = os.Getenv("OPENAI_TOKEN")
	)

	chatgpt := openai.NewChatGPT(openAIToken)

	plugins := makePlugins(chatgpt)

	loggingPlugin(plugins)

	return llmplugin.NewPluginManager(
		chatgpt,
		llmplugin.WithPlugins(plugins),
	)
}

func makePlugins(chatgpt *openai.ChatGPT) []llmplugin.Plugin {

	plugins := []llmplugin.Plugin{
		calculator.NewCalculator(),
	}

	{ // Google Search Engine
		var (
			googleEngineID = os.Getenv("GOOGLE_ENGINE_ID")
			googleToken    = os.Getenv("GOOGLE_TOKEN")
		)

		if googleEngineID != "" && googleToken != "" {
			g := google.NewGoogle(googleEngineID, googleToken, google.WithSummarizer(chatgpt))
			plugins = append(plugins, g)
		}
	}

	{ // Customize Search Engine: agi.cn search
		plugins = append(plugins, agicn_search.NewAgicnSearch())
	}

	{ // Stable Diffusion
		if sdAddr := os.Getenv("SD_ADDR"); sdAddr != "" {
			plugins = append(plugins, stablediffusion.NewStableDiffusion(sdAddr))
		}
	}

	return plugins
}

func loggingPlugin(plugins []llmplugin.Plugin) {
	for _, p := range plugins {
		logrus.Infof("load plugin: %v", p.GetName())
	}
}
