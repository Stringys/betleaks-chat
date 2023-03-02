package controller

import "github.com/gin-gonic/gin"

type ChatController struct {
}

func NewChatController() ChatController {
	return ChatController{}
}

func (ctrl ChatController) GetMembers(ctx *gin.Context) {
}

func (ctrl ChatController) GetChats(ctx *gin.Context) {
}

func (ctrl ChatController) JoinChat(ctx *gin.Context) {
}

func (ctrl ChatController) CreateChat(ctx *gin.Context) {
}
