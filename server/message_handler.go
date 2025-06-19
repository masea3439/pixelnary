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

type GameOver struct {
	Word string `json:"word"`
}

const roundTime = 90
const nextRoundWaitTime = 5
const startingGridSize = 7
const minGridSize = 3

// const wordlistKey = "p-wordlist"

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

func disconnectPlayer(conn *websocket.Conn) {
	if conn == nil {
		return
	}
	sendMessage(conn, getJsonMessage("disconnect", ""))
	conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
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
	case "rematch":
		if senderConn == gameRoom.player1Conn {
			gameRoom.player1Rematch = true
		} else {
			gameRoom.player2Rematch = true
		}
		if gameRoom.player1Rematch && gameRoom.player2Rematch {
			gameRoom.player1Rematch = false
			gameRoom.player2Rematch = false
			gameRoom.round = 0
			startNextRound(gameRoom)
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

// Storing word list in redis does not make sense as the entire word list size is minuscule
// Can be stored in memory globally instead for use by all goroutines
/**
func getRandomWord(gameRoom *GameRoom) string {
	// Using shuffle pop random method as word list is tiny thus memory use is tiny
	// Should switch to repeat random if word list becomes large
	if len(gameRoom.nextWords) == 0 {
		gameRoom.nextWords = make([]int, numWords)
		for i := 0; i < numWords; i++ {
			gameRoom.nextWords[i] = i
		}
		rand.Shuffle(numWords, func(i int, j int) {
			gameRoom.nextWords[i], gameRoom.nextWords[j] = gameRoom.nextWords[j], gameRoom.nextWords[i]
		})
	}

	var wordIndex int
	wordIndex, gameRoom.nextWords = pop(gameRoom.nextWords)

	return LIndex(wordlistKey, wordIndex)
}
*/

func setNextWord(gameRoom *GameRoom) {
	if len(gameRoom.nextWords) == 0 {
		gameRoom.nextWords = make([]string, len(fullWordList))
		copy(gameRoom.nextWords, fullWordList)
		rand.Shuffle(len(gameRoom.nextWords), func(i int, j int) {
			gameRoom.nextWords[i], gameRoom.nextWords[j] = gameRoom.nextWords[j], gameRoom.nextWords[i]
		})
	}

	gameRoom.word, gameRoom.nextWords = pop(gameRoom.nextWords)
}

func startNextRound(gameRoom *GameRoom) {
	gameRoom.round++
	gridSize := max(startingGridSize+1-gameRoom.round, minGridSize)
	drawRolePlayerId := gameRoom.round%2 + 1
	setNextWord(gameRoom)
	gameRoom.roundTimer = time.AfterFunc(time.Duration(roundTime)*time.Second, func() {
		gameOver(gameRoom)
	})
	sendStartRoundMessage(gameRoom, 1, drawRolePlayerId, gridSize)
	sendStartRoundMessage(gameRoom, 2, drawRolePlayerId, gridSize)
}

func completeRound(gameRoom *GameRoom) {
	if gameRoom.roundTimer != nil {
		gameRoom.roundTimer.Stop()
	}
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
	gameOverMessage := GameOver{gameRoom.word}
	jsonGameOverMessage, err := json.Marshal(gameOverMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonMessage := getJsonMessage("game-over", string(jsonGameOverMessage))
	sendMessage(gameRoom.player1Conn, jsonMessage)
	sendMessage(gameRoom.player2Conn, jsonMessage)
}

var fullWordList = []string{
	"Acorn",
	"Apple",
	"Arrow",
	"Axe",
	"Bagel",
	"Banana",
	"Band Aid",
	"Battery",
	"Bee",
	"Belt",
	"Bird",
	"Blood",
	"Boat",
	"Bow",
	"Branch",
	"Bridge",
	"Bubble",
	"Bucket",
	"Bus",
	"Cactus",
	"Cake",
	"Candle",
	"Candy Cane",
	"Car",
	"Carrot",
	"Castle",
	"Chair",
	"Chess",
	"Chicken",
	"City",
	"Clock",
	"Cloud",
	"Clown",
	"Coffee",
	"Coin",
	"Cookie",
	"Corn",
	"Diamond",
	"Disco Ball",
	"Door",
	"Duck",
	"Earth",
	"Eggplant",
	"Eye",
	"Fire",
	"Fish",
	"Fishing Rod",
	"Flag",
	"Flashlight",
	"Flower",
	"Fly",
	"Fork",
	"Frog",
	"Glasses",
	"Grape",
	"Grass",
	"Hamburger",
	"Hammer",
	"Hat",
	"Heart",
	"High Heel",
	"Hospital",
	"Hot Dog",
	"House",
	"Igloo",
	"Knife",
	"Ladder",
	"Leaf",
	"Light Bulb",
	"Lightbulb",
	"Lightning",
	"Magnet",
	"Moon",
	"Mountain",
	"Movie",
	"Mushroom",
	"Nest",
	"Orange",
	"Pants",
	"Pencil",
	"Pickaxe",
	"Pie",
	"Pig",
	"Pizza",
	"Plant",
	"Rainbow",
	"Ring",
	"Road",
	"Robot",
	"Rock",
	"Rocket",
	"Sand",
	"Scythe",
	"Shield",
	"Shirt",
	"Shoe",
	"Shovel",
	"Sink",
	"Skull",
	"Slide",
	"Snake",
	"Snowflake",
	"Soccer Ball",
	"Staircase",
	"Star",
	"Sun",
	"Sunflower",
	"Sword",
	"Table",
	"Tie",
	"Tire",
	"Tree",
	"Turtle",
	"Umbrella",
	"Water",
	"Watermelon",
	"Window",
	"Worm",
	"Zebra",
}
