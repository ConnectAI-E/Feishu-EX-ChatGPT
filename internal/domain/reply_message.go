package domain

import (
	"github.com/google/uuid"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type ReplyMessage larkim.ReplyMessageReq

func MakeSimpleReply(replyMessageID string, message string) *ReplyMessage {

	content := larkim.NewTextMsgBuilder().
		Text(message).Build()

	reply := larkim.NewReplyMessageReqBuilder().
		MessageId(replyMessageID).Body(
		larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			Uuid(uuid.New().String()).
			Content(content).
			Build()).
		Build()

	return (*ReplyMessage)(reply)
}
