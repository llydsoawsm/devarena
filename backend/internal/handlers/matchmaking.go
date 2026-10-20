package handlers

import (
	"net/http"

	"devcode/backend/internal/models"
	"devcode/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MatchmakingHandler struct {
	db     *gorm.DB
	engine *services.MatchmakingEngine
}

func NewMatchmakingHandler(db *gorm.DB, engine *services.MatchmakingEngine) *MatchmakingHandler {
	return &MatchmakingHandler{
		db:     db,
		engine: engine,
	}
}

type EnqueueRequest struct {
	Mode string `json:"mode" binding:"required,oneof=ranked casual tournament"`
}

func (h *MatchmakingHandler) Enqueue(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := rawUserID.(uint)

	var req EnqueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	h.engine.Enqueue(&services.QueuePlayer{
		UserID: userID,
		Rating: user.Rating,
		Mode:   req.Mode,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Enqueued successfully",
		"mode":    req.Mode,
		"rating":  user.Rating,
	})
}

func (h *MatchmakingHandler) Dequeue(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := rawUserID.(uint)

	h.engine.Dequeue(userID)
	c.JSON(http.StatusOK, gin.H{"message": "Dequeued successfully"})
}

func (h *MatchmakingHandler) Status(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := rawUserID.(uint)

	isQueued := h.engine.IsQueued(userID)
	c.JSON(http.StatusOK, gin.H{"is_queued": isQueued})
}
