package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	ws "ethosview-backend/internal/websocket"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	manager  *ws.Manager
	upgrader gorilla.Upgrader
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(manager *ws.Manager) *WebSocketHandler {
	return &WebSocketHandler{
		manager: manager,
		upgrader: gorilla.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get user ID from context if authenticated
	var userID *int
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(int); ok {
			userID = &id
		}
	}

	// Generate unique client ID
	clientID := generateClientID()

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Create new client
	client := &ws.Client{
		ID:      clientID,
		UserID:  userID,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		Manager: h.manager,
	}

	// Register client with manager
	h.manager.Register <- client

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()

	// Send welcome message
	welcomeMessage := ws.Message{
		Type: "welcome",
		Data: gin.H{
			"client_id": clientID,
			"message":   "Connected to EthosView WebSocket",
		},
	}

	if welcomeData, err := welcomeMessage.Marshal(); err == nil {
		client.Send <- welcomeData
	}
}

// GetWebSocketStatus returns WebSocket connection status
func (h *WebSocketHandler) GetWebSocketStatus(c *gin.Context) {
	clientCount := h.manager.GetClientCount()

	c.JSON(http.StatusOK, gin.H{
		"status":       "running",
		"client_count": clientCount,
		"message":      "WebSocket server is active",
	})
}

// generateClientID generates a unique client ID
func generateClientID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
