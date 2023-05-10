package domain

import (
	"context"
	"strings"
)

var helperStrings = []string{
	"/help",
	"/h",
}

// HelperAction 帮助卡片
type HelperAction struct{}

func NewHelperAction() *HelperAction {
	return &HelperAction{}
}

func (h HelperAction) Execute(ctx context.Context, actionInfo *ActionInfo) (next bool, err error) {

	text := actionInfo.GetText()
	if include := h.includeHelperString(text); !include {
		return true, nil
	}

	msgID := actionInfo.GetMessageID()

	actionInfo.Result = &ActionResult{
		ReplyMsgID: msgID,
		Type:       ActionResultHelpCard,
	}

	return false, nil
}

func (h HelperAction) includeHelperString(text string) bool {

	for _, s := range helperStrings {
		if strings.Contains(text, s) {
			return true
		}
	}

	return false
}
