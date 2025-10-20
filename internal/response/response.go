package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Status int

const (
	StatusSuccess Status = iota
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusSuccess:
		return "success"
	case StatusError:
		return "error"
	default:
		return fmt.Sprintf("unknown status (%d)", int(s))
	}
}

type ApiError string

const (
	NoError    ApiError = ""
	TestFailed ApiError = "test failed"
)

type ApiResponse struct {
	Status string   `json:"status" xml:"status"`
	Data   any      `json:"data" xml:"data"`
	Error  ApiError `json:"error,omitempty" xml:"error,omitempty"`
}

func SendResponse(c *gin.Context, httpStatus int, status Status, err ApiError, data any) {
	acceptHeader := c.GetHeader("Accept")
	if strings.Contains(acceptHeader, "application/xml") {
		c.XML(httpStatus, ApiResponse{
			Status: status.String(),
			Data:   data,
			Error:  err,
		})
	} else {
		c.JSON(httpStatus, ApiResponse{
			Status: status.String(),
			Data:   data,
			Error:  err,
		})
	}
}

func SendSuccess(c *gin.Context, data any) {
	SendResponse(c, http.StatusOK, StatusSuccess, NoError, data)
}
