package response

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, log *slog.Logger, statusCode int, message string) {
	log.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
