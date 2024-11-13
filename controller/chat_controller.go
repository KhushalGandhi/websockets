package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"websockets/service"
)

type ChatController struct {
	ChatService *service.ChatService
}

func NewChatController(chatService *service.ChatService) *ChatController {
	return &ChatController{ChatService: chatService}
}

func (cc *ChatController) CreateChatGroup(c *gin.Context) {
	groupID := uuid.New().String()
	cc.ChatService.CreateGroup(groupID)
	c.JSON(http.StatusOK, gin.H{"invite_link": "/api/chat_group/join/" + groupID})
}

func (cc *ChatController) JoinChatGroup(c *gin.Context) {
	groupID := c.Param("group_id")
	username := c.Query("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	group := cc.ChatService.GetGroup(groupID)
	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat group not found"})
		return
	}

	cc.ChatService.HandleWebSocket(c.Writer, c.Request, group, username)
}
