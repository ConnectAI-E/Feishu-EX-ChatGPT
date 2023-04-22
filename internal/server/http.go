package server

import (
	"github.com/gin-gonic/gin"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/api/router"
)

func NewHTTPServer(service router.Service) *gin.Engine {
	r := gin.Default()

	router.RegisterService(service, r)
	return r
}
