package handler

import (
	"fitnessapi/internal/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "1", "exp": time.Now().Add(24 * time.Hour).Unix()})
		signed, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": signed})
	}
}
