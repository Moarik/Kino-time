package controller

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	JwtSecret = []byte("21cb87e4067d33009b081f3a8090c16d5ae10171a6ca4ddf4b494ad55164a429")
	JwtExp    = 1440
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth header is not valid"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid format for auth"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("No signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("username", claims["username"])
		}

		c.Next()
	}
}

func AuthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("No signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}

		if userId, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userId)
		}

		c.Set("isAuthenticated", true)
		c.Next()
	}
}
