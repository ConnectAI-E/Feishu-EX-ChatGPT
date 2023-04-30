package domain

import (
	"context"
	"strings"
	"time"

	"github.com/agi-cn/llmplugin"
	"github.com/agi-cn/llmplugin/plugins/agicn_search"
	"github.com/agi-cn/llmplugin/plugins/google"
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

		// TODO(zy): Google Search result may reach the max limit of feishu message.
		if pluginCtx.GetName() == (google.Google{}).GetName() {
			result = a.getOnelineForGoogleSearch(result)
		}

		// NOTE(zy): handle search result for feishu message
		if pluginCtx.GetName() == (agicn_search.AgicnSearch{}).GetName() {
			result = a.getOnelineForAgiCnSearch(result)
		}

		results = append(results, result)
	}

	return strings.Join(results, "\n")
}

func (a PluginAction) getOnelineForGoogleSearch(result string) string {

	if len(result) == 0 {
		return ""
	}

	lines := strings.Split(result, "\n")
	if len(lines) == 0 {
		return ""
	}

	final := lines[0]

	return strings.TrimSpace(
		strings.Replace(final, "<1>", "", 1),
	)
}

func (a PluginAction) getOnelineForAgiCnSearch(result string) string {
	return strings.ReplaceAll(result, "\n", ".")
}
