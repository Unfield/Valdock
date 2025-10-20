package handler

import (
	"context"
	"log"

	"github.com/Unfield/Valdock/internal/response"
	"github.com/gin-gonic/gin"
)

type TestResponse struct {
	PingData string `json:"ping_data"`
}

func (h *Handler) TestHandler(c *gin.Context) {
	result, err := h.managementKV.Do(context.Background(), h.managementKV.B().Ping().Build()).ToString()
	if err != nil {
		log.Print(err)
		response.SendResponse(c, 400, response.StatusError, response.TestFailed, nil)
	}

	response.SendSuccess(c, TestResponse{
		PingData: result,
	})
}
