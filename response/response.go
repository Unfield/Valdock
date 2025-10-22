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
	NoError             ApiError = ""
	InternalServerError ApiError = "internal server error"
	BadRequest          ApiError = "bad request"
)

type ApiResponse struct {
	Status           string   `json:"status" xml:"status"`
	Data             any      `json:"data" xml:"data"`
	Error            ApiError `json:"error,omitempty" xml:"error,omitempty"`
	ErrorDescription string   `json:"error_description,omitempty" xml:"error_description,omitempty"`
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

func SendError(c *gin.Context, httpStatus int, err ApiError, message string) {
	acceptHeader := c.GetHeader("Accept")
	if strings.Contains(acceptHeader, "application/xml") {
		c.XML(httpStatus, ApiResponse{
			Status:           StatusError.String(),
			ErrorDescription: message,
			Error:            err,
		})
	} else {
		c.JSON(httpStatus, ApiResponse{
			Status:           StatusError.String(),
			ErrorDescription: message,
			Error:            err,
		})
	}
}
