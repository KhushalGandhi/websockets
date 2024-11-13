package main

import (
	"github.com/gin-gonic/gin"
	"websockets/config"
	"websockets/controller"
	"websockets/service"
)

func main() {
	cfg := config.LoadConfig() // Load app configurations

	router := gin.Default()
	chatService := service.NewChatService()
	chatController := controller.NewChatController(chatService)

	// Routes
	router.POST("/api/chat_group", chatController.CreateChatGroup)
	router.GET("/api/chat_group/join/:group_id", chatController.JoinChatGroup)

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.Run(cfg.ServerAddress)
}
