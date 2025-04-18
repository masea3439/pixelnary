package main

import (
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
	player1Conn    *websocket.Conn
	player2Conn    *websocket.Conn
	word           string
	round          int
	roundTimer     *time.Timer
	player1Rematch bool
	player2Rematch bool
}

var roomKeyToGameRoom = make(map[string]*GameRoom)

var allowedOrigins = []string{"https://mathieusl.com"}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Println(origin)
		return slices.Contains(allowedOrigins, origin)
	},
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

	if gameRoom.player1Conn != nil && gameRoom.player2Conn != nil {
		startNextRound(gameRoom)
	}

	for {
		_, p, err := conn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			if gameRoom.player1Conn == conn {
				disconnectPlayer(gameRoom.player2Conn)
			} else if gameRoom.player2Conn == conn {
				disconnectPlayer(gameRoom.player1Conn)
			}
			gameRoom.player1Conn = nil
			gameRoom.player2Conn = nil
			cleanupEmptyRoom(roomKey)
			return
		} else if err != nil {
			log.Println(err)
			return
		}
		ProcessClientMessage(conn, gameRoom, p)
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
		if gameRoom.roundTimer != nil {
			gameRoom.roundTimer.Stop()
		}
		delete(roomKeyToGameRoom, roomKey)
		log.Println("Deleted empty room")
	}
}

func createNewGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		roomKey := generateRandomKey(6)
		roomKeyToGameRoom[roomKey] = &GameRoom{round: 0, word: "", player1Rematch: false, player2Rematch: false}
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
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080") //TODO Use allowed origins
		w.Header().Set("Access-Control-Allow-Origin", "https://mathieusl.com:443")
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
