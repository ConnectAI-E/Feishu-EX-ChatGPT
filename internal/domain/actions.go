package domain

import (
	"context"
)

type ActionResultType string

const (
	ActionResultText     ActionResultType = "text"
	ActionResultImageB64 ActionResultType = "image_b64"
	ActionResultHelpCard ActionResultType = "help_card"
)

type ActionResult struct {
	ReplyMsgID string

	Result string // may: text, base64 for image

	Type ActionResultType
}

type ActionInfo struct {
	Message *LarkMessage

	Result *ActionResult // 返回的消息。如果已经有了返回的消息，则不执行下一阶段。
}

func (a *ActionInfo) GetMessageID() string {
	return *(a.Message.MessageId)
}

func (a *ActionInfo) GetText() string {
	if msg := a.Message; msg != nil {

		return msg.GetText()
	}

	return ""
}

func (a *ActionInfo) ExistsResult() bool {
	return a.Result != nil
}

type Action interface {
	Execute(ctx context.Context, actionInfo *ActionInfo) (next bool, err error)
}
