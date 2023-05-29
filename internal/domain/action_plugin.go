package domain

import (
	"context"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/pkg/escape"
	"github.com/agi-cn/llmplugin"
	"github.com/agi-cn/llmplugin/plugins/stablediffusion"
	"github.com/sirupsen/logrus"
)

var _ Action = (*PluginAction)(nil)

type PluginAction struct {
	manager *llmplugin.PluginManager
}

func NewPluginAction(manager *llmplugin.PluginManager) *PluginAction {

	return &PluginAction{manager: manager}
}

func (a PluginAction) Execute(ctx context.Context, actionInfo *ActionInfo) (next bool, err error) {
	if !actionInfo.UsePlugin() {
		return true, nil
	}

	text := actionInfo.GetText()
	pluginCtxs, err := a.manager.Select(context.Background(), text)
	if err != nil {
		logrus.Errorf("PluginAction Select error: %v", err)
		return false, err
	}

	result := a.makeActionResult(ctx, pluginCtxs)
	if result == nil {
		return true, nil
	}

	// NOTE(zy): 需要特别补充一下遗漏的 reply message id
	result.ReplyMsgID = actionInfo.GetMessageID()

	if result.Type != ActionResultImageB64 {
		logrus.Debugf("got plugin result: %+v", result)
	} else {
		logrus.Debugf("got plugin result in base64 image")
	}
	actionInfo.Result = result

	return true, nil
}

func (a PluginAction) makeActionResult(ctx context.Context, pluginCtxs []llmplugin.PluginContext) *ActionResult {
	// only choice one plugin
	if len(pluginCtxs) == 0 {
		return nil
	}

	pluginCtx := pluginCtxs[0]

	logrus.Debugf("run plugin=%v input=%v", pluginCtx.GetName(), pluginCtx.Input)

	result, err := pluginCtx.Do(ctx, pluginCtx.Input)
	if err != nil {
		logrus.Errorf("plugin=%s input=%v run error: %v", pluginCtx.GetName(), pluginCtx.Input, err)
		return nil
	}

	// --- 根据插件类型做多模态 ---

	// stable diffusion - answer is base64 image
	if pluginCtx.GetName() == (stablediffusion.StableDiffusion{}).GetName() {
		return &ActionResult{
			Result: result,
			Type:   ActionResultImageB64,
		}
	}

	// 普通文本
	result = escape.String(result)
	return &ActionResult{
		Result: result,
		Type:   ActionResultText,
	}
}
