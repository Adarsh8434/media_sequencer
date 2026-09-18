package handler

import (
	"mediaSequencer/backend/internal/model"
	"mediaSequencer/backend/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	Repo *repository.PlaylistRepository
}

func NewPlaylistHandler(repo *repository.PlaylistRepository) *PlaylistHandler {
	return &PlaylistHandler{Repo: repo}
}

func (h *PlaylistHandler) GetPlaylist(c *gin.Context) {
	windowID, err := strconv.Atoi(c.Param("windowId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid window id"})
		return
	}

	playlist, err := h.Repo.GetByWindowID(windowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, playlist)
}

func (h *PlaylistHandler) AddToPlaylist(c *gin.Context) {
	var playlist model.Playlist

	if err := c.ShouldBindJSON(&playlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Create(&playlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, playlist)
}
