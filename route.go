package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
		println("router:" + c.FullPath())
	})

	r.POST("/token", CreateToken)
	r.GET("/authorize", Authorize)
	return r
}
