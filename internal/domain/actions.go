package domain

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ActionInfo struct {
	Message *LarkMessage

	ReplyMsg *ReplyMessage // 返回的消息。如果已经有了返回的消息，则不执行下一阶段。
}

func (a *ActionInfo) ExistsReplyMsg() bool {
	return a.ReplyMsg != nil
}

type Action interface {
	Execute(ctx context.Context, actionInfo *ActionInfo) (next bool, err error)
}

var _ Action = &ProcessMentionAction{}

// ProcessMentionAction 是否 @ 机器人
type ProcessMentionAction struct {
	botname string
}

func NewProcessMentionAction(botname string) *ProcessMentionAction {
	return &ProcessMentionAction{botname}
}

func (a ProcessMentionAction) Execute(ctx context.Context, actionInfo *ActionInfo) (bool, error) {
	chatType := actionInfo.Message.GetChatType()

	// 私聊直接过
	if chatType.IsUserChatType() {
		return true, nil
	}

	if !chatType.IsGroupChatType() {
		logrus.Errorf("receive message chat type unsupport: %v", chatType)
		return false, nil
	}

	// 群聊判断是否提到机器人
	return actionInfo.Message.IsMentionAt(a.botname), nil
}

var _ Action = &MessageAction{}

type MessageAction struct {
	llmer LLMer
}

func NewMessageAction(llmer LLMer) *MessageAction {
	return &MessageAction{llmer: llmer}
}

func (a MessageAction) Execute(ctx context.Context, actionInfo *ActionInfo) (next bool, err error) {

	msg := actionInfo.Message.GetText()
	msgID := actionInfo.Message.ID()

	if len(msg) == 0 {
		return false, nil
	}

	messages := a.makeLlmMessages(actionInfo)
	answer, err := a.llmer.Chat(ctx, messages)
	if err != nil {
		return false, err
	}

	replyMsg := MakeSimpleReply(msgID, answer.Content)

	actionInfo.ReplyMsg = replyMsg
	return true, nil
}

func (a MessageAction) makeLlmMessages(actionInfo *ActionInfo) []LlmMessage {
	msg := actionInfo.Message.GetText()

	return []LlmMessage{
		{
			Role:    "user",
			Content: msg,
		},
	}

}
