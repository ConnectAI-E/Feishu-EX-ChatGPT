package domain

import "context"

type Feishuer interface {
	Reply(context.Context, ReplyMessage) error
}
