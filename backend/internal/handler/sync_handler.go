package handler

import (
	"mediaSequencer/backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SyncHandler struct {
	Service *service.SyncService
}

func NewSyncHandler(service *service.SyncService) *SyncHandler {
	return &SyncHandler{Service: service}
}

func (h *SyncHandler) StartSync(c *gin.Context) {
	var request struct {
		MediaID         int `json:"mediaId"`
		DurationSeconds int `json:"durationSeconds"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Service.Start(request.MediaID, request.DurationSeconds)

	c.JSON(http.StatusOK, gin.H{
		"message": "sync started",
	})
}

func (h *SyncHandler) GetSync(c *gin.Context) {
	state := h.Service.Get()

	if state == nil {
		c.JSON(http.StatusOK, gin.H{
			"active": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active":          true,
		"mediaId":         state.MediaID,
		"durationSeconds": state.DurationSeconds,
		"startedAt":       state.StartedAt,
	})
}
