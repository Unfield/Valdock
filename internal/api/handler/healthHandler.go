package handler

import (
	"time"

	"github.com/Unfield/Valdock/internal/response"
	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string `json:"status" xml:"status"`
	Timestamp string `json:"timestamp" xml:"timestamp"`
}

func (h *Handler) HealthCheckHandler(c *gin.Context) {
	response.SendSuccess(c, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
