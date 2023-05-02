package domain

import (
	"context"
	"strings"
	"time"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/pkg/escape"

	"github.com/agi-cn/llmplugin"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

var _ Action = (*ProcessUniqueAction)(nil)

type ProcessUniqueAction struct {
	processed cache.Cache
}

func NewProcessUniqueAction() *ProcessUniqueAction {
	return &ProcessUniqueAction{
		processed: *cache.New(1*time.Hour, 1*time.Hour),
	}
}

func (a *ProcessUniqueAction) Execute(ctx context.Context, actionInfo *ActionInfo) (next bool, err error) {

	msgID := actionInfo.Message.ID()

	if _, ok := a.processed.Get(msgID); ok {
		logrus.Infof("message processed: %s", actionInfo.GetContext())
		return false, nil
	}

	a.processed.Set(msgID, true, cache.DefaultExpiration)
	return true, nil
}

type ActionInfo struct {
	Message *LarkMessage

	ReplyMsg *ReplyMessage // 返回的消息。如果已经有了返回的消息，则不执行下一阶段。
}

func (a *ActionInfo) GetContext() string {
	if msg := a.Message; msg != nil {

		return msg.GetText()
	}

	return ""
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
	if actionInfo.ExistsReplyMsg() {
		return true, nil
	}

	msg := actionInfo.GetContext()
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
	replyMsg := MakeSimpleReply(msgID, result)

	actionInfo.ReplyMsg = replyMsg
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

var _ Action = (*PluginAction)(nil)

type PluginAction struct {
	manager *llmplugin.PluginManager
}

func NewPluginAction(manager *llmplugin.PluginManager) *PluginAction {

	return &PluginAction{manager: manager}
}

func (a PluginAction) Execute(ctx context.Context, actionInfo *ActionInfo) (next bool, err error) {

	text := actionInfo.Message.GetText()
	pluginCtxs, err := a.manager.Select(context.Background(), text)
	if err != nil {
		logrus.Errorf("PluginAction Select error: %v", err)
		return false, err
	}

	result := a.makeReplyMessage(ctx, pluginCtxs)

	logrus.Debugf("got plugin result: %s", result)

	result = strings.TrimSpace(result)
	if len(result) == 0 {
		return true, nil
	}

	msgID := actionInfo.Message.ID()
	replyMsg := MakeSimpleReply(msgID, result)

	actionInfo.ReplyMsg = replyMsg
	return true, nil
}

func (a PluginAction) makeReplyMessage(ctx context.Context, pluginCtxs []llmplugin.PluginContext) string {

	results := make([]string, 0, len(pluginCtxs))
	for _, pluginCtx := range pluginCtxs {
		logrus.Debugf("run plugin=%v input=%v", pluginCtx.GetName(), pluginCtx.Input)

		result, err := pluginCtx.Do(ctx, pluginCtx.Input)
		if err != nil {
			logrus.Errorf("plugin=%s input=%v run error: %v", pluginCtx.GetName(), pluginCtx.Input, err)
			continue
		}

		result = escape.String(result)
		results = append(results, result)
	}

	return strings.Join(results, "\n")
}
