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

	broadcast chan *OutboundMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Receive data from client
	receive chan ClientInBoundMessage

	gameState *game.State
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *OutboundMessage),
		receive:    make(chan ClientInBoundMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

var framesPerSecond = float32(10.0)
var millisPerFrame = time.Duration(1000.0 / framesPerSecond)
// TODO: simulation framerate should be high - but update-broadcast rate could be slower

func (h *Hub) broadcastMessage(message *OutboundMessage) {
	for client := range h.clients {
		if client.player != nil {
			select {
			case client.send <- *message:
			default:
				log.Printf("Failed to send message to client - removing client %v, %v", client.player.Name, client.player.Id)
				h.removeClient(client)
			}
		}
	}
}

func (h *Hub) removeClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		if client.player != nil {
			h.gameState.RemovePlayer(client.player)
		}
	}
}

func buildUpdatePlayersMessage(gameState *game.State) *OutboundMessage {
	var players []PlayerInfo
	for _, player := range gameState.Players {
		playerInfo := PlayerInfo{
			Id:    player.Id,
			Name:  player.Name,
			Pos:   player.Shape.Body.Position(),
			Angle: player.Shape.Body.Angle(),
		}
		players = append(players, playerInfo)
	}
	return &OutboundMessage{
		Type: UpdatePlayers,
		Data: players,
	}
}

func (h *Hub) run() {
	h.gameState = game.CreateInstance()
	ticker := time.NewTicker(time.Millisecond * millisPerFrame)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			// Evaluate a step on the game
			h.gameState.RunStep(framesPerSecond)
			h.broadcastMessage(buildUpdatePlayersMessage(h.gameState))
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			h.removeClient(client)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		case inboundMessage := <- h.receive:
			log.Printf("Received message of type %v", inboundMessage.message.Type)
			switch inboundMessage.message.Type {
			case Register:
				player := h.gameState.AddPlayer(inboundMessage.message.Data.(string))
				inboundMessage.client.player = player
				inboundMessage.client.send<-OutboundMessage{
					Type: Registered,
					Data: player.Id,
				}
			case Unregister:
				h.gameState.RemovePlayer(inboundMessage.client.player)
				inboundMessage.client.player = nil
			case RotateClockWise:
				inboundMessage.client.player.Rotate(0.1)
			case RotateCounterClockWise:
				inboundMessage.client.player.Rotate(-0.1)
			case IncreaseThrust:
				inboundMessage.client.player.AddThrust()
			}
		}
	}
}
