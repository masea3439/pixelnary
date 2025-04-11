package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	MessageType string `json:"messageType"`
	Data        string `json:"data"`
}

type Guess struct {
	Guess     string `json:"guess"`
	IsCorrect bool   `json:"isCorrect"`
}

type GameState struct {
	PlayerId         int    `json:"playerId"`
	DrawRolePlayerId int    `json:"drawRolePlayerId"`
	MatchState       string `json:"matchState"`
	RoundTimeLeft    int    `json:"roundTimeLeft"`
	GridSize         int    `json:"gridSize"`
	Word             string `json:"word"`
}

const roundTime = 90
const startingGridSize = 7

func sendMessage(conn *websocket.Conn, data []byte) {
	if conn == nil {
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println(err)
		return
	}
}

func sendMessageToPartner(senderConn *websocket.Conn, gameRoom *GameRoom, data []byte) {
	if gameRoom.player1Conn == senderConn {
		sendMessage(gameRoom.player2Conn, data)
	} else {
		sendMessage(gameRoom.player1Conn, data)
	}
}

func ProcessClientMessage(senderConn *websocket.Conn, gameRoom *GameRoom, byteMessage []byte) {
	var message Message
	err := json.Unmarshal(byteMessage, &message)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch message.MessageType {
	case "canvas":
		// TODO check game state
		// TODO input validation
		sendMessageToPartner(senderConn, gameRoom, byteMessage)
	case "guess":
		// TODO check game state
		// TODO input validation
		// TODO check correctness
		guess := Guess{message.Data, false}
		jsonGuess, err := json.Marshal(guess)
		if err != nil {
			fmt.Println(err)
			return
		}
		message.Data = string(jsonGuess)

		jsonMessage, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		sendMessage(gameRoom.player1Conn, jsonMessage)
		sendMessage(gameRoom.player2Conn, jsonMessage)
	}
}

func sendStartRoundMessage(gameRoom *GameRoom, playerId int, drawRolePlayerId int, gridSize int) {
	var gameState GameState
	if playerId == drawRolePlayerId {
		gameState = GameState{
			PlayerId:         playerId,
			DrawRolePlayerId: drawRolePlayerId,
			MatchState:       "drawing",
			RoundTimeLeft:    roundTime,
			GridSize:         gridSize,
			Word:             gameRoom.word,
		}
	} else {
		gameState = GameState{
			PlayerId:         playerId,
			DrawRolePlayerId: drawRolePlayerId,
			MatchState:       "drawing",
			RoundTimeLeft:    roundTime,
			GridSize:         gridSize,
		}
	}
	jsonGameState, err := json.Marshal(gameState)
	if err != nil {
		fmt.Println(err)
		return
	}
	message := Message{"start-round", string(jsonGameState)}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	if playerId == 1 {
		sendMessage(gameRoom.player1Conn, jsonMessage)
	} else {
		sendMessage(gameRoom.player2Conn, jsonMessage)
	}
}

// TODO keep track of past words to avoid repetition
func getRandomWord() string {
	//TODO retrieve from redis
	return "House"
}

func startNextRound(gameRoom *GameRoom) {
	if gameRoom.roundTimer != nil {
		gameRoom.roundTimer.Stop()
	}
	gameRoom.round++
	gridSize := startingGridSize + 1 - gameRoom.round
	drawRolePlayerId := gameRoom.round%2 + 1
	gameRoom.word = getRandomWord()
	gameRoom.roundTimer = time.AfterFunc(time.Duration(roundTime)*time.Second, func() {
		gameOver(gameRoom)
	})
	sendStartRoundMessage(gameRoom, 1, drawRolePlayerId, gridSize)
	sendStartRoundMessage(gameRoom, 2, drawRolePlayerId, gridSize)
}

func gameOver(gameRoom *GameRoom) {

}
