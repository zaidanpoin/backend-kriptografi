package Middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse dan validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return privateKey, nil
		})

		// Kalau token invalid atau error (expired, format salah, dsb.)
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Ambil claim (opsional, kalau butuh user_id, email, dll)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["id"])
			// Bisa tambah c.Set("email", claims["email"]) dsb.
		}

		// Token valid, lanjutkan request
		c.Next()
	}
}
