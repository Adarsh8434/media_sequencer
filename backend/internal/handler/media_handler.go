package handler

import (
	"mediaSequencer/backend/internal/model"
	"mediaSequencer/backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	Repo *repository.MediaRepository
}

func NewMediaHandler(repo *repository.MediaRepository) *MediaHandler {
	return &MediaHandler{Repo: repo}
}

func (h *MediaHandler) GetMedia(c *gin.Context) {
	media, err := h.Repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, media)
}

func (h *MediaHandler) CreateMedia(c *gin.Context) {
	var media model.Media

	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Create(&media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, media)
}
