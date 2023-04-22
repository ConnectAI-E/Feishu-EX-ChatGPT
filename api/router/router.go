package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ConnectAI-E/Feishu-EX-ChatGPT/api/protos"
)

type Service interface {
	Hello(context.Context, *protos.HelloReq) (*protos.HelloResp, error)
}

func RegisterService(service Service, r gin.IRouter) {

	r.GET("/hello", makeHello(service))
}

func makeHello(s Service) func(c *gin.Context) {

	return func(c *gin.Context) {

		var req protos.HelloReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"result": false,
			})
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"result": true,
		})

	}
}
