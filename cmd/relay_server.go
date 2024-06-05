package main

import (
    "log"
    "net/http"
    "os"
    "os/signal"
    "github.com/gorilla/websocket"
)

// Define a struct to represent the Nostr relay server.
type NostrRelayServer struct {
	clients       map[string]*websocket.Conn
	addClient     chan *websocket.Conn
	removeClient  chan *websocket.Conn
	broadcast     chan []byte
	upgrader      websocket.Upgrader
}

// Initialize the Nostr relay server.
func NewNostrRelayServer() *NostrRelayServer {
	return &NostrRelayServer{
		clients:       make(map[string]*websocket.Conn),
		addClient:     make(chan *websocket.Conn),
		removeClient:  make(chan *websocket.Conn),
		broadcast:     make(chan []byte),
		upgrader:      websocket.Upgrader{},
	}
}

// Start the Nostr relay server.
func (s *NostrRelayServer) Start() {
	// Start a goroutine to handle incoming client connections.
	go func() {
		for {
			select {
			case conn := <-s.addClient:
				s.clients[conn.RemoteAddr().String()] = conn
			case conn := <-s.removeClient:
				delete(s.clients, conn.RemoteAddr().String())
				conn.Close()
			case msg := <-s.broadcast:
				// Broadcast the received message to all connected clients.
				for _, conn := range s.clients {
					conn.WriteMessage(websocket.TextMessage, msg)
				}
			}
		}
	}()

	// Start a goroutine to handle HTTP requests to upgrade to WebSocket connections.
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}
		defer conn.Close()

		// Add the client to the relay server.
		s.addClient <- conn
		defer func() {
			s.removeClient <- conn
		}()

		// Listen for messages from the client.
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// Broadcast the received message to all connected clients.
			s.broadcast <- msg
		}
	})

	// Start the HTTP server.
	go func() {
		log.Println("Starting Nostr Relay Server on port 8080...")
		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Fatal("Error starting Nostr Relay Server:", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C) to gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down Nostr Relay Server...")
}

func main() {
	// Create a new instance of the Nostr relay server.
	server := NewNostrRelayServer()

	// Start the Nostr relay server.
	server.Start()
}
