package http

import (
	"github.com/gin-gonic/gin"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// Error represents an API error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success sends a successful JSON response
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
	})
}

// Error sends an error JSON response
func SendError(c *gin.Context, status int, code string, message string) {
	c.JSON(status, Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}

// Paginated sends a paginated JSON response
func Paginated(c *gin.Context, status int, data interface{}, pagination interface{}) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
		Meta:    pagination,
	})
}
