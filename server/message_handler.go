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

// TODO keep track of past words to avoid repetition
func getRandomWord() string {
	//TODO retrieve from redis
	words := []string{
		"acorn", "alien", "apple", "bee", "cactus", "chair", "diamond", "duck",
		"door", "earth", "fire", "fish", "hat", "house", "key", "leaf", "ladybug",
		"mountain", "moon", "owl", "pencil", "rainbow", "snake", "skull", "soccer",
		"spider", "star", "tent", "turtle", "water", "zebra", "Strawberry", "Eclipse", "Chandelier", "Ketchup", "Toothpaste", "Rainbow", "Bunk bed", "Boardgame", "Beehive", "Lemon", "Wreath", "Waffles", "Bubble", "Whistle", "Snowball", "Bouquet", "Headphones", "Fireworks", "Igloo", "Ferris wheel", "Banana peel", "Lawnmower", "Summer", "Whisk", "Cupcake", "Sleeping bag", "Bruise", "Fog", "Crust", "Battery",
		"Giraffe",
		"Koala",
		"Wasp",
		"Scorpion",
		"Lion",
		"Salamander",
		"Dolphin",
		"Frog",
		"Panda",
		"Platypus",
		"Meerkat",
		"Eagle",
		"Buckle",
		"Lipstick",
		"Raindrop",
		"Bus",
		"Lobster",
		"Robot",
		"Car accident",
		"Lollipop",
		"Salamander",
		"Castle",
		"Magnet",
		"Slipper",
		"Chainsaw",
		"Megaphone",
		"Snowball",
		"Circus tent",
		"Mermaid",
		"Sprinkler",
		"Computer",
		"Minivan",
		"Stapler",
		"Desk",
		"Pay cheque",
		"Work computer",
		"Fax machine",
		"Phone",
		"Paper",
		"Light",
		"Chair",
		"Desk lamp",
		"Notepad",
		"Paper clips",
		"Binder",
		"Calculator",
		"Calendar",
		"Sticky Notes",
		"Pens",
		"Pencils",
		"Notebook",
		"Book",
		"Chairs",
		"Coffee cup",
		"Chairs",
		"Coffee mug",
		"Thermos",
		"Hot cup",
		"Glue",
		"Clipboard",
		"Paperclips",
		"Chocolate",
		"Secretary",
		"Work",
		"Paperwork",
		"Workload",
		"Employee",
		"Boredom",
		"Coffee",
		"Golf",
		"Laptop",
		"Sandcastle",
		"Monday",
		"Vanilla",
		"Bamboo",
		"Sneeze",
		"Scratch",
		"Celery",
		"Hammer",
		"Frog",
		"Tennis",
		"Hot dog",
		"Pants",
		"Bridge",
		"Bubblegum",
		"Candy bar",
		"Bucket",
		"Skiing",
		"Sledding",
		"Snowboarding",
		"Snowman",
		"Polar bear",
		"Cream",
		"Waffle",
		"Pancakes",
		"Ice cream",
		"Sundae",
		"beach",
		"Sunglasses",
		"Surfboard",
		"Watermelon",
		"Baseball",
		"Bat",
		"Ball",
		"T-shirt",
		"Kiss",
		"Jellyfish",
		"Jelly",
		"Butterfly",
		"Spider",
		"Broom",
		"Spiderweb",
		"Mummy",
		"Candy",
		"Bays",
		"Squirrels",
		"Basketball",
		"Water Bottle",
		"Unicorn",
		"Dog leash",
		"Newspaper",
		"Hammock",
		"Video camera",
		"Money",
		"Smiley face",
		"Umbrella",
		"Picnic basket",
		"Teddy bear",
		"Ambulance",
		"Ancient Pyramids",
		"Bacteria",
		"Goosebumps",
		"Pizza",
		"Platypus",
		"Clam Chowder",
		"Goldfish bowl",
		"Skull",
		"Spiderweb",
		"Smoke",
		"Tree",
		"Ice",
		"Blanket",
		"Seaweed",
		"Flame",
		"Bubble",
		"Hair",
		"Tooth",
		"Leaf",
		"Worm",
		"Sky",
		"Apple",
		"Plane",
		"Cow",
		"House",
		"Dog",
		"Car",
		"Bed",
		"Furniture",
		"Train",
		"Rainbow",
		"Paintings",
		"Drawing",
		"Cup",
		"Plate",
		"Bowl",
		"Cushion",
		"Sofa",
		"Sheet",
		"Kitchen",
		"Table",
		"Candle",
		"Shirt",
		"Clothes",
		"Dress",
		"Pillow",
		"Home",
		"Toothpaste",
		"Guitar",
		"Schoolbag",
		"Pencil Case",
		"Glasses",
		"Towel",
		"Watch",
		"Piano",
		"Pen",
		"Hat",
		"Shoes",
		"Socks",
		"Jeans",
		"Hair Gel",
		"Keyboard",
		"Bra",
		"Jacket",
		"Tie",
		"Bandage",
		"Scarf",
		"Hair Brush",
		"Cell Phone"}
	return words[rand.Intn(len(words))]
}

func startNextRound(gameRoom *GameRoom) {
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
