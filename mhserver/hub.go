package main

import (
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan message
	process   chan *Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan message),
		process:    make(chan *Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.broadcast:
			//	log.Println("broadcasting")
			for client := range h.clients {
				//		log.Printf("found %s\n", client.device.Name)
				for _, targ := range msg.Targets {
					if client.device.Type == targ || client.device.Name == targ {
						select {
						case client.send <- msg:
						default:
							log.Printf("Hub closing client %s\n", client.device.Name)
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		}
	}
}
