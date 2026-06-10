package middleware

import (
	"fitnessapi/internal/constants"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		logger.Error("request failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": constants.ErrInvalidRequest, "message": err.Error()})
	}
}
