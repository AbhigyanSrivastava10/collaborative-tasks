package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	clients    map[string]map[*Client]bool // boardID -> set of clients
	mu         sync.RWMutex
	redis      *redis.Client
	subscribe  chan *Client
	unsubscribe chan *Client
	broadcast  chan Message
}

type Message struct {
	BoardID string `json:"board_id"`
	Type    string `json:"type"`    // "task_created" | "task_updated" | "task_deleted"
	Payload any    `json:"payload"`
}

func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		clients:     make(map[string]map[*Client]bool),
		redis:       redisClient,
		subscribe:   make(chan *Client),
		unsubscribe: make(chan *Client),
		broadcast:   make(chan Message, 256),
	}
}

func (h *Hub) Run(ctx context.Context) {
	// Subscribe to Redis channel for broadcasts
	pubsub := h.redis.Subscribe(ctx, "tasks")
	defer pubsub.Close()

	go func() {
		for msg := range pubsub.Channel() {
			var m Message
			if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				log.Println("ws: failed to unmarshal message:", err)
				continue
			}
			h.sendToBoard(m)
		}
	}()

	for {
		select {
		case client := <-h.subscribe:
			h.mu.Lock()
			if h.clients[client.boardID] == nil {
				h.clients[client.boardID] = make(map[*Client]bool)
			}
			h.clients[client.boardID][client] = true
			h.mu.Unlock()
			fmt.Printf("ws: client joined board %s\n", client.boardID)

		case client := <-h.unsubscribe:
			h.mu.Lock()
			if clients, ok := h.clients[client.boardID]; ok {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.clients, client.boardID)
				}
			}
			h.mu.Unlock()
			close(client.send)
			fmt.Printf("ws: client left board %s\n", client.boardID)

		case msg := <-h.broadcast:
			// Publish to Redis so all server instances get it
			data, _ := json.Marshal(msg)
			h.redis.Publish(ctx, "tasks", data)

		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) sendToBoard(msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for client := range h.clients[msg.BoardID] {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients[msg.BoardID], client)
		}
	}
}

func (h *Hub) Broadcast(msg Message) {
	h.broadcast <- msg
}
