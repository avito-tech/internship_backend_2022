package formatter

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidData      = errors.New("validation error")
	ErrNotFound         = errors.New("not found error")
	ErrPermissionDenied = errors.New("permission denied")
)

type Error struct {
	Message   string `json:"message" binding:"required"`
	RequestID string `json:"request_id,omitempty"`
	Code      int    `json:"code"`
}

func NewError(message string, code int) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}

func HandleErrors(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()
	if err == nil {
		return
	}

	code := http.StatusInternalServerError
	body := NewError(err.Error(), code)

	if errors.Is(err, ErrInvalidData) {
		code = http.StatusBadRequest
		body = nil
	}

	if errors.Is(err, ErrNotFound) {
		code = http.StatusNotFound
		body = nil
	}

	if errors.Is(err, ErrPermissionDenied) {
		code = http.StatusForbidden
		body = nil
	}

	if code == http.StatusInternalServerError {
		c.Header("Retry-After", "30")
	}

	if body == nil {
		c.JSON(code, gin.H{})
	} else {
		c.JSON(code, body)
	}
}
