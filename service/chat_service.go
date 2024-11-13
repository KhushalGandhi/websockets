package service

import (
	"sync"
	"websockets/model"
)

type ChatService struct {
	Groups map[string]*model.ChatGroup
	mu     sync.RWMutex
}

func NewChatService() *ChatService {
	return &ChatService{
		Groups: make(map[string]*model.ChatGroup),
	}
}

func (cs *ChatService) CreateGroup(id string) *model.ChatGroup {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	group := &model.ChatGroup{
		ID:         id,
		Clients:    make(map[*model.Client]bool),
		Broadcast:  make(chan model.Message),
		Register:   make(chan *model.Client),
		Unregister: make(chan *model.Client),
	}
	cs.Groups[id] = group
	go cs.runGroup(group)
	return group
}

func (cs *ChatService) GetGroup(id string) *model.ChatGroup {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.Groups[id]
}

func (cs *ChatService) runGroup(group *model.ChatGroup) {
	for {
		select {
		case client := <-group.Register:
			group.Clients[client] = true
			cs.notify(group, model.Message{Username: "System", Content: client.Username + " joined the chat"})
		case client := <-group.Unregister:
			delete(group.Clients, client)
			close(client.Send)
			cs.notify(group, model.Message{Username: "System", Content: client.Username + " left the chat"})
		case msg := <-group.Broadcast:
			for client := range group.Clients {
				client.Send <- msg
			}
		}
	}
}

func (cs *ChatService) notify(group *model.ChatGroup, msg model.Message) {
	group.Broadcast <- msg
}
