package domain

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type FeishuEx struct {
	actions  []Action
	feishuer Feishuer
}

type FeishuExOption func(*FeishuEx)

// WithAction 添加单个处理流程
func WithAction(action Action) FeishuExOption {
	return func(fe *FeishuEx) {
		fe.actions = append(fe.actions, action)
	}
}

// WithActions 添加多个处理流程
func WithActions(actions ...Action) FeishuExOption {
	return func(fe *FeishuEx) {
		fe.actions = append(fe.actions, actions...)
	}
}

func NewFeishuEx(feishuer Feishuer, opts ...FeishuExOption) (*FeishuEx, error) {
	feishuEx := &FeishuEx{
		feishuer: feishuer,
		actions:  make([]Action, 0, 8),
	}

	for _, opt := range opts {
		opt(feishuEx)
	}

	if err := feishuEx.checkInitSucc(); err != nil {
		return nil, err
	}

	return feishuEx, nil
}

func (f *FeishuEx) checkInitSucc() error {
	if len(f.actions) == 0 {
		return errors.New("empty actions")
	}

	return nil
}

func (f *FeishuEx) HandleMessageReceive(ctx context.Context, msg *LarkMessage) error {
	actionInfo := f.makeActionInfo(msg)

	if err := f.processActionInfo(ctx, actionInfo); err != nil {
		logrus.Errorf("process action_info error: %v", err)
		return err
	}

	return f.sendReplyMessage(ctx, actionInfo.Result)
}

func (f *FeishuEx) makeActionInfo(msg *LarkMessage) *ActionInfo {
	return &ActionInfo{
		Message: msg,
	}
}

func (f *FeishuEx) processActionInfo(ctx context.Context, actionInfo *ActionInfo) error {

	for _, action := range f.actions {
		next, err := action.Execute(ctx, actionInfo)
		if err != nil {
			logrus.Errorf("action execute error: %v", err)
			return err
		}
		if !next {
			return nil
		}

		if actionInfo.ExistsResult() {
			return nil
		}
	}

	return nil
}

func (f *FeishuEx) sendReplyMessage(ctx context.Context, actionResult *ActionResult) error {
	// NOTE(zy): 这里有可能是没有结果，也有可能是重复消息处理。
	if actionResult == nil {
		logrus.Warnf("sendReplyMessage: reply failed, nil reply message")
		return nil
	}

	return f.feishuer.Reply(ctx, actionResult)
}
