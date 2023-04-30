package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"result": true,
		})
	})

	return r
}
