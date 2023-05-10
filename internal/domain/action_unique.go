package domain

import (
	"context"
	"time"

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
		logrus.Infof("message processed: %s", actionInfo.GetText())
		return false, nil
	}

	a.processed.Set(msgID, true, cache.DefaultExpiration)
	return true, nil
}
