package service

import (
	"context"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/internal/domain"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type FeishuExService struct {
	feishuEx *domain.FeishuEx
}

func NewFeishuExService(feishuEx *domain.FeishuEx) *FeishuExService {
	return &FeishuExService{
		feishuEx: feishuEx,
	}
}

func (s *FeishuExService) HandleMessageReceive(ctx context.Context, receive *larkim.P2MessageReceiveV1) error {

	larkMessage := (*domain.LarkMessage)(receive.Event.Message)

	return s.feishuEx.HandleMessageReceive(ctx, larkMessage)
}
