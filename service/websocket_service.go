package service

import (
	"github.com/gorilla/websocket"
	"net/http"
	"websockets/model"
	"websockets/utils"
)

func (cs *ChatService) HandleWebSocket(w http.ResponseWriter, r *http.Request, group *model.ChatGroup, username string) {
	conn, err := utils.UpgradeToWebSocket(w, r)
	if err != nil {
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}

	client := &model.Client{
		Username: username,
		Group:    group, // Assign the group here
		Send:     make(chan model.Message),
	}
	group.Register <- client

	go cs.writeMessages(client, conn)
	cs.readMessages(client, conn)
}

func (cs *ChatService) readMessages(client *model.Client, conn *websocket.Conn) {
	defer func() {
		client.Group.Unregister <- client // Use client.Group here
		conn.Close()
	}()
	for {
		var msg model.Message
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		msg.Username = client.Username
		client.Group.Broadcast <- msg // Use client.Group here
	}
}

func (cs *ChatService) writeMessages(client *model.Client, conn *websocket.Conn) {
	for msg := range client.Send {
		if err := conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
