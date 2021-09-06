package ws

import (
	"encoding/json"
)

type (
	// Hub maintains the set of active clients and broadcasts messages to the
	// clients.
	Hub struct {
		// Registered clients.
		clients map[*Client]bool

		// Inbound messages from the clients.
		broadcast chan Broadcast

		// Register requests from the clients.
		register chan *Client

		// Unregister requests from clients.
		unregister chan *Client

		subscribe chan subscription

		unsubscribe chan subscription

		channels map[string]map[*Client]bool
	}

	WS interface {
		Broadcast(message Broadcast)
	}

	subscription struct {
		client  *Client
		channel string
	}
)

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan Broadcast),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		subscribe:   make(chan subscription),
		unsubscribe: make(chan subscription),
		channels:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			m, _ := json.Marshal(message)
			for client := range h.channels[message.Channel] {
				select {
				case client.send <- m:
				default:
					h.unregisterClient(client)
				}
			}
		case message := <-h.subscribe:
			switch message.channel {
			case BlocksChannel, TransactionsChannel:
			default:
				continue
			}
			if _, ok := h.channels[message.channel]; !ok {
				h.channels[message.channel] = make(map[*Client]bool)
			}
			h.channels[message.channel][message.client] = true
		case message := <-h.unsubscribe:
			for client := range h.channels[message.channel] {
				delete(h.channels[message.channel], client)
			}
		}
	}
}

func (h *Hub) Broadcast(message Broadcast) {
	h.broadcast <- message
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		for _, clients := range h.channels {
			delete(clients, client)
		}
		close(client.send)
	}
}
