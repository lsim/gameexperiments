package main

import (
	"log"
	"github.com/lsim/gameexperiments/backend/game"
	"time"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan OutBoundMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Receive data from client
	receive chan ClientInBoundMessage
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan OutBoundMessage),
		receive:    make(chan ClientInBoundMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

var framesPerSecond = 20.0
var millisPerFrame = time.Duration(1000.0 / framesPerSecond)

func (h *Hub) broadcastMessage(message OutBoundMessage) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) removeClient(client *Client, gameState *game.State) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		gameState.RemovePlayer(client.player)
	}
}

func (h *Hub) run() {
	gameState := game.CreateInstance()
	ticker := time.NewTicker(time.Millisecond * millisPerFrame)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			// Evaluate a step on the game
			gameState.RunStep()
			h.broadcastMessage(OutBoundMessage{*gameState})
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			h.removeClient(client, gameState)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		case inboundMessage := <- h.receive:
			log.Printf("Received message of type %v", inboundMessage.message.Type)
			switch inboundMessage.message.Type {
			case Register:
				player := gameState.AddPlayer(inboundMessage.message.Data.(string))
				inboundMessage.client.player = player
			}
		}
	}
}
