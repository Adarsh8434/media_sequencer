package handler

import (
	"mediaSequencer/backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WindowHandler struct {
	Repo *repository.WindowRepository
}

func NewWindowHandler(repo *repository.WindowRepository) *WindowHandler {
	return &WindowHandler{Repo: repo}
}

func (h *WindowHandler) GetWindows(c *gin.Context) {
	windows, err := h.Repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, windows)
}
