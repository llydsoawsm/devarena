package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type EventMessage struct {
	Event   string      `json:"event"`
	Channel string      `json:"channel,omitempty"`
	Payload interface{} `json:"payload"`
}

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   uint
	MatchID  string // active public match ID
	IsClosed bool
	mu       sync.Mutex
}

type Room struct {
	MatchPublicID string
	Clients       map[*Client]bool
	mu            sync.RWMutex
}

type Hub struct {
	Clients    map[*Client]bool
	Rooms      map[string]*Room
	UserMap    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan EventMessage
	mu         sync.RWMutex
}

var GlobalHub *Hub

func InitHub() *Hub {
	h := &Hub{
		Clients:    make(map[*Client]bool),
		Rooms:      make(map[string]*Room),
		UserMap:    make(map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan EventMessage, 256),
	}
	GlobalHub = h
	go h.Run()
	return h
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			if client.UserID > 0 {
				h.UserMap[client.UserID] = client
			}
			h.mu.Unlock()
			log.Printf("WS: Client connected (UserID: %d)", client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				if client.UserID > 0 {
					delete(h.UserMap, client.UserID)
				}
				if client.MatchID != "" {
					if room, exists := h.Rooms[client.MatchID]; exists {
						room.mu.Lock()
						delete(room.Clients, client)
						if len(room.Clients) == 0 {
							delete(h.Rooms, client.MatchID)
						}
						room.mu.Unlock()
					}
				}
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("WS: Client disconnected (UserID: %d)", client.UserID)

		case msg := <-h.Broadcast:
			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			h.mu.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- data:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) JoinRoom(client *Client, matchID string) {
	h.mu.Lock()
	room, exists := h.Rooms[matchID]
	if !exists {
		room = &Room{
			MatchPublicID: matchID,
			Clients:       make(map[*Client]bool),
		}
		h.Rooms[matchID] = room
	}
	client.MatchID = matchID
	room.mu.Lock()
	room.Clients[client] = true
	room.mu.Unlock()
	h.mu.Unlock()
}

func (h *Hub) BroadcastToRoom(matchID string, event string, payload interface{}) {
	h.mu.RLock()
	room, exists := h.Rooms[matchID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	msg := EventMessage{
		Event:   event,
		Channel: "match." + matchID,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()
	for client := range room.Clients {
		select {
		case client.Send <- data:
		default:
		}
	}
}

func (h *Hub) SendToUser(userID uint, event string, payload interface{}) {
	h.mu.RLock()
	client, exists := h.UserMap[userID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	msg := EventMessage{
		Event:   event,
		Channel: "user",
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
	default:
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var eventMsg EventMessage
		if err := json.Unmarshal(message, &eventMsg); err != nil {
			continue
		}

		switch eventMsg.Event {
		case "join_room":
			if matchID, ok := eventMsg.Payload.(string); ok {
				c.Hub.JoinRoom(c, matchID)
			}
		case "participant.state":
			if c.MatchID != "" {
				c.Hub.BroadcastToRoom(c.MatchID, "participant.state", eventMsg.Payload)
			}
		case "participant.heartbeat":
			// Acknowledge heartbeat
			c.Send <- []byte(`{"event":"pong"}`)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(25 * time.Second)
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
