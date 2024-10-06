package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// MARK: Types
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

// MARK: Functions

/**
 * Creates and returns a new Hub instance.
 * The Hub is responsible for managing WebSocket clients
 * and orchestrating message flow.
 */
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// MARK: Methods

/**
 * Starts the main loop of the Hub.
 * This method handles client registration, unregistration,
 * and message distribution.
 */
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			h.mutex.Lock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}

/**
 * Broadcasts a message to all connected clients.
 * The message is first marshaled to JSON before broadcasting.
 */
func (h *Hub) Broadcast(message Message) error {
	json, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}
	h.broadcast <- json
	return nil
}

/**
 * ReadPump handles reading messages from a client's WebSocket connection.
 * It continuously pulls data from the client into the Hub.
 */
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("error unmarshaling message: %v", err)
			continue
		}

		// Handle the message based on its type
		switch msg.Type {
		case "join_game":
			// Handle join game logic
		case "player_action":
			// Handle player action logic
		default:
			fmt.Printf("unknown message type: %s", msg.Type)
		}
	}
}

/**
 * WritePump handles writing messages to a client's WebSocket connection.
 * It continuously pushes data from the Hub out to the client.
 */
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		message, ok := <-c.Send
		if !ok {
			// The hub closed the channel.
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}
