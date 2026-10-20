package handlers

import (
	"net/http"

	"devcode/backend/internal/models"
	"devcode/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MatchHandler struct {
	db     *gorm.DB
	engine *services.MatchmakingEngine
}

func NewMatchHandler(db *gorm.DB, engine *services.MatchmakingEngine) *MatchHandler {
	return &MatchHandler{
		db:     db,
		engine: engine,
	}
}

func (h *MatchHandler) GetByPublicID(c *gin.Context) {
	publicID := c.Param("public_id")

	var match models.Match
	if err := h.db.
		Preload("Problem").
		Preload("Problem.TestCases", "is_hidden = ?", false).
		Preload("Participants.User").
		Where("public_id = ?", publicID).
		First(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"match": match})
}

func (h *MatchHandler) Forfeit(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := rawUserID.(uint)
	publicID := c.Param("public_id")

	if err := h.engine.ForfeitMatch(publicID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Forfeited successfully"})
}

func (h *MatchHandler) Leaderboard(c *gin.Context) {
	var users []models.User
	if err := h.db.Select("id, name, rating, matches_played, wins, losses").
		Order("rating DESC").
		Limit(50).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": users})
}

func (h *MatchHandler) History(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := rawUserID.(uint)

	var participants []models.MatchParticipant
	if err := h.db.Preload("User").
		Where("user_id = ?", userID).
		Order("id DESC").
		Limit(20).
		Find(&participants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch match history"})
		return
	}

	matchIDs := make([]uint, len(participants))
	for i, p := range participants {
		matchIDs[i] = p.MatchID
	}

	var matches []models.Match
	if len(matchIDs) > 0 {
		h.db.Preload("Problem").
			Preload("Participants.User").
			Where("id IN ?", matchIDs).
			Order("id DESC").
			Find(&matches)
	}

	c.JSON(http.StatusOK, gin.H{"matches": matches})
}
