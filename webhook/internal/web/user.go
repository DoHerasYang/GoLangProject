package main

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")
	userGroup.POST("/signup", h.SignUp)
	userGroup.POST("/login", h.Login)
	userGroup.GET("/profile", h.Profile)
	userGroup.POST("/edit", h.Edit)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {

}

func (h *UserHandler) Login(ctx *gin.Context) {

}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func (h *UserHandler) Profile(ctx *gin.Context) {

}
