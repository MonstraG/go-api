// Copyright 2013 The Gorilla WebSocket Authors.
// This code has been since modified by me.

package websockets

import (
	"go-server/setup/reqRes"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte
}

func newHub() *Hub {
	return &Hub{
		broadcast: make(chan []byte),
		clients:   make(map[*Client]bool),
	}
}

var HubSingleton = newHub()

func (hub *Hub) RegisterClient(client *Client) {
	hub.clients[client] = true
}

func (hub *Hub) UnRegisterClient(client *Client) {
	if _, ok := hub.clients[client]; ok {
		delete(hub.clients, client)
		close(client.send)
	}
}

func (hub *Hub) Broadcast(message string) {
	for client := range hub.clients {
		select {
		case client.send <- []byte(message):
		default:
			close(client.send)
			delete(hub.clients, client)
		}
	}
}

// ServeWs handles websocket requests from the peer.
func (hub *Hub) ServeWs(w reqRes.MyWriter, r *reqRes.MyRequest) {
	conn, err := upgrader.Upgrade(w, &r.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	hub.RegisterClient(client)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
