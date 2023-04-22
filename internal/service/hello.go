package service

import (
	"context"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/api/protos"
	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/api/router"
)

var _ router.Service = &ExChatGPTService{}

type ExChatGPTService struct {
}

func NewExChatGPTService() *ExChatGPTService {
	return &ExChatGPTService{}
}

func (*ExChatGPTService) Hello(ctx context.Context, req *protos.HelloReq) (*protos.HelloResp, error) {
	name := req.Name

	return &protos.HelloResp{
		Message: "hello, " + name,
	}, nil
}
