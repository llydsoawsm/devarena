package ws

import (
	"log"
	"net/http"

	"devcode/backend/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

func ServeWS(hub *Hub, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		var userID uint = 0

		if tokenStr != "" {
			if claims, err := auth.ValidateToken(tokenStr, jwtSecret); err == nil {
				userID = claims.UserID
			}
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}

		client := &Client{
			Hub:     hub,
			Conn:    conn,
			Send:    make(chan []byte, 256),
			UserID:  userID,
			MatchID: c.Query("match_id"),
		}

		hub.Register <- client

		if client.MatchID != "" {
			hub.JoinRoom(client, client.MatchID)
		}

		go client.WritePump()
		go client.ReadPump()
	}
}
