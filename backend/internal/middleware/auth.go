package middleware

import (
	"fitnessapi/internal/config"
	"fitnessapi/internal/constants"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Next()
			return
		}
		tokenText := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenText, func(token *jwt.Token) (interface{}, error) { return []byte(cfg.JWTSecret), nil })
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": constants.ErrUnauthorized, "message": "无效或过期的 JWT"})
			return
		}
		c.Next()
	}
}
