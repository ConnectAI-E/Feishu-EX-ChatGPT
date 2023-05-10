package feishu

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/domain"
	"github.com/google/uuid"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var _ domain.Feishuer = &Feishu{}

type Feishu struct {
	client *lark.Client
}

func NewFeishu(client *lark.Client) *Feishu {
	return &Feishu{client}
}

func (f Feishu) Reply(ctx context.Context, actionResult *domain.ActionResult) error {

	switch actionResult.Type {
	case domain.ActionResultText:
		return f.replyText(ctx, actionResult)

	case domain.ActionResultHelpCard:
		return f.replyHelpCard(ctx, actionResult)

	case domain.ActionResultImageB64:
		return f.replyImageByBase64(ctx, actionResult)

	}

	return errors.New("unknown type for reply message")
}

func (f Feishu) sendReplyMessage(ctx context.Context, reply *larkim.ReplyMessageReq) error {

	resp, err := f.client.Im.Message.Reply(ctx, reply)
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

func (f Feishu) replyText(ctx context.Context, actionResult *domain.ActionResult) error {
	msg := actionResult.Result

	msg = strings.TrimSpace(msg)
	if len(msg) == 0 {
		return nil
	}

	reply := MakeSimpleReply(actionResult.ReplyMsgID, msg)

	return f.sendReplyMessage(ctx, reply)
}

func (f Feishu) replyHelpCard(ctx context.Context, actionResult *domain.ActionResult) error {

	reply := MakeSendHelpCard(actionResult.ReplyMsgID)

	return f.sendReplyMessage(ctx, reply)
}

func (f Feishu) replyImageByBase64(ctx context.Context, actionResult *domain.ActionResult) error {
	var (
		msgID  = actionResult.ReplyMsgID
		b64Img = actionResult.Result
	)

	payload, err := base64.StdEncoding.DecodeString(b64Img)
	if err != nil {
		return err
	}

	imageKey, err := f.uploadImagePayload(ctx, payload)
	if err != nil {
		return err
	}

	reply, err := f.makeImageReplyMessage(msgID, imageKey)
	if err != nil {
		return err
	}

	return f.sendReplyMessage(ctx, reply)
}

func (f Feishu) uploadImagePayload(ctx context.Context, payload []byte) (string, error) {

	resp, err := f.client.Im.Image.Create(context.Background(),
		larkim.NewCreateImageReqBuilder().
			Body(larkim.NewCreateImageReqBodyBuilder().
				ImageType(larkim.ImageTypeMessage).
				Image(bytes.NewReader(payload)).
				Build()).
			Build())
	if err != nil {
		return "", errors.Wrap(err, "upload image to feishu")
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return "", errors.New(resp.Msg)
	}
	return *(resp.Data.ImageKey), nil
}

func (f Feishu) makeImageReplyMessage(msgID string, imageKey string) (*larkim.ReplyMessageReq, error) {
	msgImage := larkim.MessageImage{ImageKey: imageKey}
	msgImageContent, err := msgImage.String()
	if err != nil {
		return nil, err
	}

	reply := larkim.NewReplyMessageReqBuilder().
		MessageId(msgID).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeImage).
			Uuid(uuid.New().String()).
			Content(msgImageContent).
			Build()).
		Build()

	return reply, nil
}
