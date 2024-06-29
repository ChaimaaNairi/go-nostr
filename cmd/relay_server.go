// relay_server.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message represents the structure of the message sent over WebSocket
type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow any origin for WebSocket connections
		},
	}
)

func main() {
	fmt.Println("Starting relay server...")

	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// Implement message routing logic here (forwarding messages, etc.)
		log.Printf("Received message: %s from %s to %s\n", message.Content, message.Sender, message.Recipient)
	}
}
