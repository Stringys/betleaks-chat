package main

import (
	"log"

	"github.com/Stringys/betleaks-chat/config"
	"github.com/Stringys/betleaks-chat/controller"
	"github.com/Stringys/betleaks-chat/db"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Load()
	_, err := db.NewPostgresDB(config.DBConStr)
	if err != nil {
		log.Fatalf("failed to connect PostgresDB: %v", err)
	}

	chatController := controller.NewChatController()

	engine := gin.Default()
	engine.GET("/members", chatController.GetMembers)
	engine.GET("/chats", chatController.GetChats)
	engine.GET("/join", chatController.JoinChat)
	engine.GET("/create-chat", chatController.CreateChat)
	engine.Run()

}
