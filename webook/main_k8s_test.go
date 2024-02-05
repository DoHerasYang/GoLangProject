package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello!")
	})
	server.Run(":8080")
}
