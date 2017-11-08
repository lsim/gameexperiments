package main

import (
	"log"
	"github.com/lsim/gameexperiments/backend/game"
	"time"
)

var (
	framesPerSecond = float32(100.0)
	broadcastsPerSecond = float32(10.0)
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

func (h *Hub) broadcastMessage(message *OutboundMessage) {
	for client := range h.clients {
		select {
		case client.send <- *message:
		default:
			log.Printf("Failed to send message to client - removing client %v, %v", client.player.Name, client.player.Id)
			h.removeClient(client)
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

func (h *Hub) resetClientPlayerRef(player *game.Player) {
	for client := range h.clients {
		if client.player == player {
			client.player = nil
			break
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
			Velocity: player.Shape.Body.Velocity(),
		}
		players = append(players, playerInfo)
	}
	return &OutboundMessage{
		Type: UpdatePlayers,
		Data: players,
	}
}

func buildPlayerDeadMessage(player *game.Player) *OutboundMessage {
	return &OutboundMessage{
		Type: PlayerDied,
		Data: player.Id,
	}
}

func buildBulletDeadMessage(bullet *game.PlayerBullet) *OutboundMessage {
	return &OutboundMessage{
		Type: BulletDied,
		Data: bullet.Id,
	}
}

func buildShootMessage(bullet *game.PlayerBullet) *OutboundMessage {
	body := bullet.Shape.Body
	return &OutboundMessage{
		Type: Shoot,
		Data: BulletInfo{
			Id: bullet.Id,
			Pos: body.Position(),
			Angle: body.Angle(),
			Velocity: body.Velocity(),
		},
	}
}

func buildWelcomeMessage() *OutboundMessage {
	return &OutboundMessage{
		Type: WelcomeClient,
		Data: WelcomeMessage{
			RotateFactor: game.RotateFactor,
			PlanetRadius: game.PlanetRadius,
			GravityStrength: game.GravityStrength,
			PlayerLength: game.PlayerLength,
			PlayerMass: game.PlayerMass,
			ThrustFactor: game.ThrustFactor,
		},
	}
}

func (h *Hub) run() {
	h.gameState = game.CreateInstance()
	frameTicker := time.NewTicker(time.Millisecond * time.Duration(1000.0 / framesPerSecond))
	broadcastTicker := time.NewTicker(time.Millisecond * time.Duration(1000.0 / broadcastsPerSecond))
	defer func() {
		frameTicker.Stop()
		broadcastTicker.Stop()
	}()

	for {
		select {
		case <- broadcastTicker.C:
			h.broadcastMessage(buildUpdatePlayersMessage(h.gameState))
		case <-frameTicker.C:
			h.gameState.RunStep(framesPerSecond)
		case deadPlayer := <-h.gameState.PlayerDeaths:
			h.broadcastMessage(buildPlayerDeadMessage(deadPlayer))
			h.resetClientPlayerRef(deadPlayer)
			h.gameState.RemovePlayer(deadPlayer)
		case deadBullet := <-h.gameState.BulletDeaths:
			h.gameState.RemoveBullet(deadBullet)
			h.broadcastMessage(buildBulletDeadMessage(deadBullet))
		case client := <-h.register:
			h.clients[client] = true
			client.send<-*buildWelcomeMessage()
		case client := <-h.unregister:
			deadPlayer := client.player
			h.removeClient(client)
			// For now just treat quitting as dying to get some cleanup at other clients
			if deadPlayer != nil {
			  h.broadcastMessage(buildPlayerDeadMessage(deadPlayer))
			}
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		case inboundMessage := <- h.receive:
			//log.Printf("Received message of type %v", inboundMessage.message.Type)
			switch inboundMessage.message.Type {
			case Register:
				player := h.gameState.AddPlayer(inboundMessage.message.Data.(string))
				inboundMessage.client.player = player
				inboundMessage.client.send<-OutboundMessage{
					Type: Registered,
					Data: player.Id,
				}
			case Unregister:
				if player := inboundMessage.client.player; player != nil {
					h.gameState.RemovePlayer(player)
					h.broadcastMessage(buildPlayerDeadMessage(player))
				}
				inboundMessage.client.player = nil
			case RotateClockWise:
				inboundMessage.client.player.Rotate(1)
			case RotateCounterClockWise:
				inboundMessage.client.player.Rotate(-1)
			case IncreaseThrust:
				inboundMessage.client.player.AddThrust()
			case Shoot:
				if bullet, ok := inboundMessage.client.player.Shoot(); ok {
					h.broadcastMessage(buildShootMessage(bullet))
				}
			}
		}
	}
}
