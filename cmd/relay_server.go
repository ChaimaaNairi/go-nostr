package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
	upgrader  = websocket.Upgrader{}
	messages  []Message
	mu        sync.Mutex
)

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/events", handleEvents)

	go handleMessages()

	log.Println("Relay server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		messages = append(messages, msg)
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
