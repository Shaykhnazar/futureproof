package handlers

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"

	"github.com/shaykhnazar/futureproof/internal/models"
)

// WebSocketHub manages WebSocket connections
type WebSocketHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan models.WebSocketMessage
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.RWMutex
	logger     *zap.Logger
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub(logger *zap.Logger) *WebSocketHub {
	hub := &WebSocketHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan models.WebSocketMessage, 256),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		logger:     logger,
	}

	go hub.run()
	return hub
}

// run manages the WebSocket hub lifecycle
func (h *WebSocketHub) run() {
	for {
		select {
		case conn := <-h.register:
			h.mutex.Lock()
			h.clients[conn] = true
			h.mutex.Unlock()
			h.logger.Info("WebSocket client connected", zap.Int("total_clients", len(h.clients)))

		case conn := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mutex.Unlock()
			h.logger.Info("WebSocket client disconnected", zap.Int("total_clients", len(h.clients)))

		case message := <-h.broadcast:
			h.mutex.RLock()
			data, _ := json.Marshal(message)
			for conn := range h.clients {
				err := conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					h.logger.Error("Failed to send message to client", zap.Error(err))
					go func(c *websocket.Conn) {
						h.unregister <- c
					}(conn)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *WebSocketHub) Broadcast(messageType string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to marshal broadcast data", zap.Error(err))
		return
	}

	message := models.WebSocketMessage{
		Type:      messageType,
		Timestamp: time.Now(),
		Data:      jsonData,
	}

	select {
	case h.broadcast <- message:
	default:
		h.logger.Warn("Broadcast channel full, message dropped")
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHub) HandleWebSocket(c *websocket.Conn) {
	// Register client
	h.register <- c

	// Send welcome message
	welcome := models.WebSocketMessage{
		Type:      "connected",
		Timestamp: time.Now(),
		Data:      json.RawMessage(`{"message":"Connected to FutureProof real-time updates"}`),
	}
	data, _ := json.Marshal(welcome)
	c.WriteMessage(websocket.TextMessage, data)

	// Handle messages
	defer func() {
		h.unregister <- c
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			h.logger.Debug("WebSocket connection closed", zap.Error(err))
			break
		}

		// Echo message back (or handle specific client messages)
		h.logger.Debug("Received WebSocket message", zap.ByteString("message", message))
	}
}

// UpgradeMiddleware upgrades HTTP connection to WebSocket
func UpgradeMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
