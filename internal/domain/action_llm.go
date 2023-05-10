package domain

import (
	"context"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/pkg/escape"
	"github.com/sirupsen/logrus"
)

var _ Action = &MessageAction{}

type MessageAction struct {
	llmer LLMer
}

func NewMessageAction(llmer LLMer) *MessageAction {
	return &MessageAction{llmer: llmer}
}

func (a MessageAction) Execute(ctx context.Context, actionInfo *ActionInfo) (next bool, err error) {
	if actionInfo.ExistsResult() {
		return true, nil
	}

	msg := actionInfo.GetText()
	msgID := actionInfo.Message.ID()

	if len(msg) == 0 {
		return false, nil
	}

	messages := a.makeLlmMessages(msg)
	answer, err := a.llmer.Chat(ctx, messages)
	if err != nil {
		logrus.Errorf("MessageAction: llmer chat error: %v", err)
		return false, err
	}

	logrus.Debugf("MessageAction: llmer.Chat: message=%s answer: %s", msg, answer.Content)

	result := escape.String(answer.Content)

	actionInfo.Result = &ActionResult{
		ReplyMsgID: msgID,
		Result:     result,
		Type:       ActionResultText,
	}
	return true, nil
}

func (a MessageAction) makeLlmMessages(content string) []LlmMessage {
	return []LlmMessage{
		{
			Role:    "user",
			Content: content,
		},
	}

}
