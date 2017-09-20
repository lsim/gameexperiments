package main

import (
	"net/http"
	"encoding/json"
	"log"
	//"github.com/gorilla/websocket"
	//"time"
	//"./game"
	//"github.com/stojg/vector"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Message string
	}{
		"Hello, World",
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
}

//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//}

//type SocketMessage struct {
//	Message string
//	Counter int
//	*game.State
//}

//func socketHandler(w http.ResponseWriter, r *http.Request) {
//	//conn, err := upgrader.Upgrade(w, r, nil)
//	//if err != nil {
//	//	log.Println(err)
//	//	return
//	//}
//
//	//for i := 0; ; i++ {
//	//	time.Sleep(time.Second)
//	//	conn.WriteJSON(SocketMessage{Message:"Message " + string(i), Counter: i})
//	//}
//	// Add a player to the game and associate the player with the conn
//	// When data is received on the conn, we need to update the player object
//	// When a round has been simulated, we need to broadcast the state to all conns
//	// Looks like we need a map from conn -> player
//	//conn.ReadJSON()
//}

//type serverState struct {
//	gameState game.State
//	playerByConn map[*websocket.Conn]*game.Player
//}

//var state = serverState{}
//
//func (state *serverState)registerPlayer(conn *websocket.Conn, name string) {
//	player := state.gameState.AddPlayer(name)
//	state.playerByConn[conn] = player
//}
//
func main() {
	// Spin off the hub
	hub := newHub()
	go hub.run()

	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("./node_modules/"))))
	http.Handle("/", http.FileServer(http.Dir("./frontend/")))

	http.HandleFunc("/api/socket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/api/hello", helloWorld)
	port := ":4567"
	log.Println("Server listening at localhost" + port)
	http.ListenAndServe(port, nil)

}
