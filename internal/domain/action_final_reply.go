package domain

import "context"

var _ Action = (*FinalReply)(nil)

type FinalReply struct {
	text string
}

func NewFinalReply(text string) *FinalReply {
	return &FinalReply{text}
}

func (a *FinalReply) Execute(ctx context.Context, actionInfo *ActionInfo) (next bool, err error) {

	if actionInfo.Result == nil {
		actionInfo.Result = &ActionResult{
			ReplyMsgID: actionInfo.GetMessageID(),
			Type:       ActionResultText,
			Result:     a.text,
		}
	}

	return true, nil
}
