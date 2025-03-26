package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var allowedOrigins = []string{"ws://localhost:3333", "wss://localhost:3333"}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
		// origin := r.Header.Get("Origin")
		// return slices.Contains(allowedOrigins, origin)
	},
}

func openWebSocketConn(w http.ResponseWriter, r *http.Request) {
	log.Println("Upgrading")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
		log.Println(p)
	}
}

func getHello(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello")
}

func main() {
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/ws", openWebSocketConn)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
