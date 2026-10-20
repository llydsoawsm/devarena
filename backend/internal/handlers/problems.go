package handlers

import (
	"net/http"

	"devcode/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProblemHandler struct {
	db *gorm.DB
}

func NewProblemHandler(db *gorm.DB) *ProblemHandler {
	return &ProblemHandler{db: db}
}

func (h *ProblemHandler) List(c *gin.Context) {
	var problems []models.Problem
	if err := h.db.Select("id, slug, title, difficulty, time_limit_ms, memory_limit_kb, created_at").Find(&problems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch problems"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"problems": problems})
}

func (h *ProblemHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var problem models.Problem
	if err := h.db.Preload("TestCases", "is_hidden = ?", false).Where("slug = ?", slug).First(&problem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"problem": problem})
}

func (h *ProblemHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	var problem models.Problem
	if err := h.db.Preload("TestCases", "is_hidden = ?", false).First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"problem": problem})
}
