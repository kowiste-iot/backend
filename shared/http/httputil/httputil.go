package httputil

import (
	"backend/shared/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func NewErrorResponse(c *gin.Context, err error) {
	var appErr *errors.AppError
	if e, ok := err.(*errors.AppError); ok {
		appErr = e
	} else {
		appErr = errors.NewInternal("Internal server error", err)
	}

	c.JSON(appErr.Code, Response{
		Success: false,
		Error: gin.H{
			"type":    appErr.Type,
			"message": appErr.Message,
		},
	})
}
