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
