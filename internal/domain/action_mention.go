package domain

import (
	"context"

	"github.com/sirupsen/logrus"
)

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
