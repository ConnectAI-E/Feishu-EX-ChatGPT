package main

import (
	"os"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/domain"
	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/repo/feishu"
	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/server"
	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/service"

	"github.com/joho/godotenv"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("load env failed: %v", err)
	}

	var (
		feishuAppID       = os.Getenv("FEISHU_APP_ID")
		feishuAppSecret   = os.Getenv("FEISHU_APP_SECRET")
		botname           = os.Getenv("BOT_NAME")
		feishuVerifyToken = os.Getenv("VERIFY_TOKEN")
		feishuEncryptKey  = os.Getenv("ENCRYPT_KEY")

		port = os.Getenv("HTTP_PORT")
	)

	feishuClient := lark.NewClient(feishuAppID, feishuAppSecret)
	feishuer := feishu.NewFeishu(feishuClient)

	feishuEx, err := domain.New(
		feishuer,
		domain.WithActions(
			domain.NewProcessMentionAction(botname),
			domain.NewMessageAction(),
		),
	)
	if err != nil {
		logrus.Fatalf("FeishuEx init failed: %v", err)
	}

	service := service.NewFeishuExService(feishuEx)

	feishuDispatcher := dispatcher.NewEventDispatcher(feishuVerifyToken, feishuEncryptKey).
		OnP2MessageReceiveV1(service.HandleMessageReceive)

	r := server.NewHTTPServer(service)

	r.POST("/webhook/event", sdkginext.NewEventHandlerFunc(feishuDispatcher))
	logrus.Fatal(r.Run(port))
}
