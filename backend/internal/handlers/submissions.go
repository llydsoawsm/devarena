package handlers

import (
	"net/http"
	"time"

	"devcode/backend/internal/models"
	"devcode/backend/internal/services"
	"devcode/backend/internal/ws"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubmissionHandler struct {
	db     *gorm.DB
	judge  *services.JudgeService
	engine *services.MatchmakingEngine
	hub    *ws.Hub
}

func NewSubmissionHandler(db *gorm.DB, judge *services.JudgeService, engine *services.MatchmakingEngine, hub *ws.Hub) *SubmissionHandler {
	return &SubmissionHandler{
		db:     db,
		judge:  judge,
		engine: engine,
		hub:    hub,
	}
}

type RunSampleRequest struct {
	ProblemID uint   `json:"problem_id" binding:"required"`
	Language  string `json:"language" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

type SubmitMatchRequest struct {
	MatchPublicID string `json:"match_public_id" binding:"required"`
	Language      string `json:"language" binding:"required"`
	Code          string `json:"code" binding:"required"`
}

func (h *SubmissionHandler) RunSample(c *gin.Context) {
	var req RunSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var testCases []models.TestCase
	if err := h.db.Where("problem_id = ? AND is_hidden = ?", req.ProblemID, false).
		Order("sequence_order ASC").
		Find(&testCases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load sample test cases"})
		return
	}

	res := h.judge.Execute(req.Language, req.Code, testCases)
	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (h *SubmissionHandler) SubmitMatch(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := rawUserID.(uint)

	var req SubmitMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var match models.Match
	if err := h.db.Preload("Participants.User").Where("public_id = ?", req.MatchPublicID).First(&match).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if match.Status != "live" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match is not live"})
		return
	}

	var participant *models.MatchParticipant
	var opponent *models.MatchParticipant
	for i := range match.Participants {
		if match.Participants[i].UserID == userID {
			participant = &match.Participants[i]
		} else {
			opponent = &match.Participants[i]
		}
	}

	if participant == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not a participant in this match"})
		return
	}

	var testCases []models.TestCase
	if err := h.db.Where("problem_id = ?", match.ProblemID).Order("sequence_order ASC").Find(&testCases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load test cases"})
		return
	}

	res := h.judge.Execute(req.Language, req.Code, testCases)

	// Record submission
	submission := models.Submission{
		MatchID:         &match.ID,
		UserID:          userID,
		ProblemID:       match.ProblemID,
		Language:        req.Language,
		Code:            req.Code,
		Verdict:         res.Verdict,
		PassedTestCases: res.PassedTestCases,
		TotalTestCases:  res.TotalTestCases,
		RuntimeMs:       res.RuntimeMs,
		MemoryKb:        res.MemoryKb,
		CompileOutput:   res.ErrorDetails,
		SubmittedAt:     time.Now(),
	}
	h.db.Create(&submission)

	// Update participant stats
	participant.SubmissionCount++
	participant.PassedTests = res.PassedTestCases
	participant.TotalTests = res.TotalTestCases
	h.db.Save(participant)

	// Broadcast opponent progress to room
	h.hub.BroadcastToRoom(match.PublicID, "participant.state", map[string]interface{}{
		"user_id":      userID,
		"passed_tests": res.PassedTestCases,
		"total_tests":  res.TotalTestCases,
		"verdict":      res.Verdict,
	})

	// Broadcast submission result notification
	h.hub.BroadcastToRoom(match.PublicID, "match.submission_result", map[string]interface{}{
		"user_id":      userID,
		"verdict":      res.Verdict,
		"passed_tests": res.PassedTestCases,
		"total_tests":  res.TotalTestCases,
	})

	// If all passed (AC), match is won!
	if res.Verdict == "AC" && opponent != nil {
		go h.engine.SettleMatch(&match, participant, opponent, false)
	}

	c.JSON(http.StatusOK, gin.H{
		"submission": submission,
		"result":     res,
	})
}
