package ws

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // in production, check the origin
	},
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusBadRequest)
		return
	}

	client := NewClient(boardID, conn, h.hub)
	h.hub.subscribe <- client

	go client.WritePump()
	go client.ReadPump()
}
