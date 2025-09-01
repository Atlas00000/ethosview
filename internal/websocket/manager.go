package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Message represents a WebSocket message
type Message struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	Time   time.Time   `json:"time"`
	UserID *int        `json:"user_id,omitempty"`
}

// Marshal marshals the message to JSON
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Client represents a WebSocket client connection
type Client struct {
	ID      string
	UserID  *int
	Conn    *websocket.Conn
	Send    chan []byte
	Manager *Manager
	mu      sync.Mutex
}

// Manager handles WebSocket connections and message broadcasting
type Manager struct {
	clients    map[string]*Client
	broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

// NewManager creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Start starts the WebSocket manager
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			m.mu.Lock()
			m.clients[client.ID] = client
			m.mu.Unlock()
			log.Printf("Client %s connected", client.ID)

		case client := <-m.Unregister:
			m.mu.Lock()
			if _, ok := m.clients[client.ID]; ok {
				delete(m.clients, client.ID)
				close(client.Send)
			}
			m.mu.Unlock()
			log.Printf("Client %s disconnected", client.ID)

		case message := <-m.broadcast:
			m.mu.RLock()
			for _, client := range m.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.clients, client.ID)
				}
			}
			m.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (m *Manager) Broadcast(messageType string, data interface{}) {
	message := Message{
		Type: messageType,
		Data: data,
		Time: time.Now(),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	m.broadcast <- jsonData
}

// BroadcastToUser sends a message to a specific user
func (m *Manager) BroadcastToUser(userID int, messageType string, data interface{}) {
	message := Message{
		Type:   messageType,
		Data:   data,
		Time:   time.Now(),
		UserID: &userID,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	m.mu.RLock()
	for _, client := range m.clients {
		if client.UserID != nil && *client.UserID == userID {
			select {
			case client.Send <- jsonData:
			default:
				close(client.Send)
				delete(m.clients, client.ID)
			}
		}
	}
	m.mu.RUnlock()
}

// GetClientCount returns the number of connected clients
func (m *Manager) GetClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

// WritePump handles writing messages to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
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
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump handles reading messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.Manager.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Handle incoming messages (e.g., subscription requests)
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Handle different message types
		switch msg.Type {
		case "ping":
			response := Message{
				Type: "pong",
				Time: time.Now(),
			}
			responseData, _ := json.Marshal(response)
			c.Send <- responseData
		case "subscribe":
			// Handle subscription requests
			log.Printf("Client %s subscribed to updates", c.ID)
		}
	}
}
