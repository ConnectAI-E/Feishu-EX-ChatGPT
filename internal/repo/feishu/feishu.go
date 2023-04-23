package feishu

import (
	"context"
	"errors"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/domain"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
)

var _ domain.Feishuer = &Feishu{}

type Feishu struct {
	client *lark.Client
}

func NewFeishu(client *lark.Client) *Feishu {
	return &Feishu{client}
}

func (f Feishu) Reply(ctx context.Context, replyMsg domain.ReplyMessage) error {

	resp, err := f.client.Im.Message.Reply(ctx, (*larkim.ReplyMessageReq)(&replyMsg))
	if err != nil {
		logrus.Errorf("feishu reply failed: %v", err)
		return err
	}

	if !resp.Success() {
		logrus.Errorf("feishu reply not success: %v", resp)
		return errors.New("feishu reply not success")
	}

	return nil
}
