package server

import (
	"github.com/gin-gonic/gin"
)

func NewHTTPServer() *gin.Engine {
	r := gin.Default()

	return r
}
