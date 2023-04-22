package main

import (
	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/domain/feishuex"
	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/server"
	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/service"

	sdkginext "github.com/larksuite/oapi-sdk-gin"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/sirupsen/logrus"
)

func main() {

	feishuEx := feishuex.New()
	service := service.NewExChatGPTService()

	feishuDispatcher := dispatcher.NewEventDispatcher("v", "e").
		OnP2MessageReceiveV1(feishuEx.HandleMessageReceive)

	r := server.NewHTTPServer(service)

	r.POST("/webhook/event", sdkginext.NewEventHandlerFunc(feishuDispatcher))

	logrus.Fatal(r.Run(":8080"))
}
