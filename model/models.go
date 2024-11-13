package model

type ChatGroup struct {
	ID         string
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Username string
	Group    *ChatGroup // Reference to the chat group this client belongs to
	Send     chan Message
}

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}
