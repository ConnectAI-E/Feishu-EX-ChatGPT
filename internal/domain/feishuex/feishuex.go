package feishuex

import (
	"context"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
)

type FeishuEx struct {
}

func New() *FeishuEx {
	return &FeishuEx{}
}

func (FeishuEx) HandleMessageReceive(ctx context.Context, event *larkim.P2MessageReceiveV1) error {

	logrus.Infof("receive event: %s", event.Body)

	return nil
}
