package main

import (
	"net/http"
	"log"
)

// Fire up this main func and run 'npm run watchify' on the command line to get the show on the road
// We then get full browserify + hot module reload + vue.js (compiled templates) + phaser.io
// Oh joy! :D
func main() {
	// Spin off the hub
	hub := newHub()
	go hub.run()

	http.Handle("/frontend/dist/", http.StripPrefix("/frontend/dist/", http.FileServer(http.Dir("./frontend/dist/"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	// Serve index.html specifically
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/api/socket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	port := ":4567"
	log.Println("Server listening at localhost" + port)
	http.ListenAndServe(port, nil)

}
