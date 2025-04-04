package main

import (
	"encoding/json"
	"fmt"
	"log"

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
		var guess = Guess{message.Data, false}
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
