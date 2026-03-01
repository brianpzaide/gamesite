package internal

import (
	"errors"
	"sync"

	"hunaidsav/gamesite/games"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

func (h *Hub) CreateRoom(constructor func() games.Game, gameType string) string {
	id := uuid.New().String()

	room := NewRoom(id, constructor(), gameType, h)

	h.mu.Lock()
	h.rooms[id] = room
	h.mu.Unlock()

	go room.Run()

	return id
}

func (h *Hub) RoomExists(id string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[id]
	if !ok {
		return "", false
	}

	return room.gameType, true
}

func (h *Hub) AddClient(roomID string, conn *websocket.Conn) error {
	h.mu.RLock()
	room, ok := h.rooms[roomID]
	h.mu.RUnlock()

	if !ok {
		return errors.New("room not found")
	}

	client := NewClient(conn, room)
	room.register <- client

	return nil
}

func (h *Hub) Shutdown() {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, room := range h.rooms {
		close(room.stop)
		delete(h.rooms, room.ID)
	}
}
