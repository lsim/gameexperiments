package main

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/websocket"
	"time"
)

// helloWorld handler returns a json response with a 'Hello, World' message.
//
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SocketMessage struct {
	Message string
	Counter int
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; ; i++ {
		time.Sleep(time.Second)
		conn.WriteJSON(SocketMessage{"Message " + string(i), i})
	}

	//... Use conn to send and receive messages.
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./frontend/")))

	http.HandleFunc("/api/socket", socketHandler)
	http.HandleFunc("/api/hello", helloWorld)
	port := ":4567"
	http.ListenAndServe(port, nil)
	log.Println("Server listening at localhost" + port)
}
