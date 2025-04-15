package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	MessageType string `json:"messageType"`
	Data        string `json:"data"`
}

type CompletedMessage struct {
	TimeUntilNextRound int `json:"timeUntilNextRound"`
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
const nextRoundWaitTime = 10
const startingGridSize = 7
const minGridSize = 3

func getJsonMessage(messageType string, data string) []byte {
	message := Message{messageType, data}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return jsonMessage
}

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
		isCorrect := strings.EqualFold(message.Data, gameRoom.word)

		guess := Guess{message.Data, isCorrect}

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

		if isCorrect {
			completeRound(gameRoom)
		}
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
	jsonMessage := getJsonMessage("start-round", string(jsonGameState))
	if playerId == 1 {
		sendMessage(gameRoom.player1Conn, jsonMessage)
	} else {
		sendMessage(gameRoom.player2Conn, jsonMessage)
	}
}

// TODO keep track of past words to avoid repetition
func getRandomWord() string {
	//TODO retrieve from redis
	words := []string{
		"acorn", "alien", "apple", "bee", "cactus", "chair", "diamond", "duck",
		"door", "earth", "fire", "fish", "hat", "house", "key", "leaf", "ladybug",
		"mountain", "moon", "owl", "pencil", "rainbow", "snake", "skull", "soccer",
		"spider", "star", "tent", "turtle", "water", "zebra"}
	return words[rand.Intn(len(words))]
}

func startNextRound(gameRoom *GameRoom) {
	if gameRoom.roundTimer != nil {
		gameRoom.roundTimer.Stop()
	}
	gameRoom.round++
	gridSize := max(startingGridSize+1-gameRoom.round, minGridSize)
	drawRolePlayerId := gameRoom.round%2 + 1
	gameRoom.word = getRandomWord()
	gameRoom.roundTimer = time.AfterFunc(time.Duration(roundTime)*time.Second, func() {
		gameOver(gameRoom)
	})
	sendStartRoundMessage(gameRoom, 1, drawRolePlayerId, gridSize)
	sendStartRoundMessage(gameRoom, 2, drawRolePlayerId, gridSize)
}

func completeRound(gameRoom *GameRoom) {
	completedMessage := CompletedMessage{nextRoundWaitTime}
	jsonCompletedMessage, err := json.Marshal(completedMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonMessage := getJsonMessage("round-completed", string(jsonCompletedMessage))
	sendMessage(gameRoom.player1Conn, jsonMessage)
	sendMessage(gameRoom.player2Conn, jsonMessage)
	time.AfterFunc(nextRoundWaitTime*time.Second, func() { startNextRound(gameRoom) }) //TODO check if room is still open
}

func gameOver(gameRoom *GameRoom) {

}
