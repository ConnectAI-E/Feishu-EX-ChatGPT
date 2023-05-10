package main

import (
	"context"
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

	plugins := makePlugins()

	loggingPlugin(plugins)

	return llmplugin.NewPluginManager(
		chatgpt,
		llmplugin.WithPlugins(plugins),
	)
}

func makePlugins() []llmplugin.Plugin {

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

	{ // Google Search Engine
		var (
			googleEngineID = os.Getenv("GOOGLE_ENGINE_ID")
			googleToken    = os.Getenv("GOOGLE_TOKEN")
		)

		if googleEngineID != "" && googleToken != "" {
			plugins = append(plugins,
				google.NewGoogle(googleEngineID, googleToken))
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
