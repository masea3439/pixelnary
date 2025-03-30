package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/gorilla/websocket"
)

type GameRoom struct {
	player1Conn *websocket.Conn
	player2Conn *websocket.Conn
}

var roomKeyToGameRoom = make(map[string]*GameRoom)

var allowedOrigins = []string{"http://localhost:8080", "https://localhost:8080"}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Println(origin)
		return slices.Contains(allowedOrigins, origin)
	},
}

func sendMessageToPartner(conn *websocket.Conn, gameRoom *GameRoom, data []byte) {
	if gameRoom.player1Conn == conn {
		gameRoom.player2Conn.WriteMessage(websocket.TextMessage, data)
	} else {
		gameRoom.player1Conn.WriteMessage(websocket.TextMessage, data)
	}
}

func openWebSocketConn(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	roomKey := queryParams.Get("key")

	gameRoom, exists := roomKeyToGameRoom[roomKey]
	if !exists {
		http.Error(w, "Room does not exist", http.StatusBadRequest)
		log.Println("Room does not exist")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if gameRoom.player1Conn == nil {
		gameRoom.player1Conn = conn
	} else if gameRoom.player2Conn == nil {
		gameRoom.player2Conn = conn
	} else {
		conn.WriteMessage(websocket.PingMessage, []byte("full")) //TODO
		conn.Close()
		log.Println("Room full")
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

		//TODO Input validation
		//TODO Save to database
		//TODO Verify game state
		canvas := bytes.Split(p, []byte(","))
		log.Println(canvas)
		sendMessageToPartner(conn, gameRoom, p)
	}
}

func generateRandomKey(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func cleanupEmptyRoom(roomKey string) {
	gameRoom, exists := roomKeyToGameRoom[roomKey]
	if exists && gameRoom.player1Conn == nil && gameRoom.player2Conn == nil {
		delete(roomKeyToGameRoom, roomKey)
		log.Println("Deleted empty room")
	}
}

func createNewGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		roomKey := generateRandomKey(6)
		roomKeyToGameRoom[roomKey] = &GameRoom{}
		// delete room after timeout if no websocket connection made
		time.AfterFunc(5*time.Second, func() { cleanupEmptyRoom(roomKey) })
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(roomKey))
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080") //TODO Use allowed origins
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		next.ServeHTTP(w, r)
	}
}

func main() {
	http.HandleFunc("/ws", middleware(openWebSocketConn))
	http.HandleFunc("/api/host", middleware(createNewGame))

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
