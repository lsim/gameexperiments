package main

import (
	"github.com/gorilla/websocket"
	"time"
	"log"
	"net/http"
	"github.com/lsim/gameexperiments/backend/game"
	"github.com/vova616/chipmunk/vect"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type OutboundMessage struct {
	Type MessageType
	Data interface{}
}

type PlayerInfo struct {
	Id    int
	Name  string
	Pos   vect.Vect
	Angle vect.Float
	Velocity vect.Vect
}

type BulletInfo struct {
	Id int
	Pos vect.Vect
	Angle vect.Float
	Velocity vect.Vect
}

type WelcomeMessage struct {
	PlayerMass        float32
	PlayerLength      float32
	PlanetRadius      float32
	GravityStrength   vect.Float
	ThrustFactor      int
	RotateFactor      float32
}

type MessageType int

const (
	Register MessageType = iota
	UpdatePlayers MessageType = iota
	Registered MessageType = iota
	Unregister MessageType = iota
	RotateClockWise MessageType = iota
	RotateCounterClockWise MessageType = iota
	IncreaseThrust MessageType = iota
	PlayerDied MessageType = iota
	Shoot MessageType = iota
	BulletDied MessageType = iota
	WelcomeClient MessageType = iota
)

type InBoundMessage struct {
	Type MessageType
	Data interface{}
}

type ClientInBoundMessage struct {
	message InBoundMessage
	client *Client
}

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan OutboundMessage

	receive chan InBoundMessage

	player *game.Player
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		message := InBoundMessage{}
		if err := c.conn.ReadJSON(&message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.receive <- ClientInBoundMessage{message: message, client: c}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:     hub,
		conn:    conn,
		send:    make(chan OutboundMessage, 64),
		receive: make(chan InBoundMessage, 64),
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
