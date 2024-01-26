package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func DefaultRouter(engine *gin.Engine) *gin.Engine {
	r := gin.Default()
	r.GET("/hello", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Hello Gin",
		})
	})
	return r
}
