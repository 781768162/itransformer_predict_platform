package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gateway/pkg/jwt"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "missing Authorization header"})
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "invalid Authorization header format"})
            c.Abort()
            return
        }

        claims, err := jwt.ParseToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "invalid or expired token"})
            c.Abort()
            return
        }

        c.Set("userId", claims.UserID)
        c.Set("userName", claims.UserName)

        c.Next()
	}
}
