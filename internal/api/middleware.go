package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
			return
		}

		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		}, jwt.WithValidMethods([]string{"HS256"}))

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}

		c.Next()
	}
}
