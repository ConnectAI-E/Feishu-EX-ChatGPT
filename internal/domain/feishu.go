package domain

import "context"

type Feishuer interface {
	Reply(ctx context.Context, actionResult *ActionResult) error
}
