package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Функция для генерации JWT токена
func generateJWT(secret []byte, valid bool) string {
	claims := jwt.MapClaims{
		"username": "test_user",
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
	}
	
	var token *jwt.Token
	if valid {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"exp": time.Now().Unix() - 1000})
	}
	
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func TestJWTMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid token passes", func(t *testing.T) {
		r := gin.New()
		r.Use(JWTMiddleware())
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
		})

		token := generateJWT(JwtSecret, true)
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing Authorization header", func(t *testing.T) {
		r := gin.New()
		r.Use(JWTMiddleware())
		r.GET("/protected", func(c *gin.Context) {})

		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Authorization format", func(t *testing.T) {
		r := gin.New()
		r.Use(JWTMiddleware())
		r.GET("/protected", func(c *gin.Context) {})

		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "WrongFormatToken")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid token", func(t *testing.T) {
		r := gin.New()
		r.Use(JWTMiddleware())
		r.GET("/protected", func(c *gin.Context) {})

		token := generateJWT(JwtSecret, false) // Токен просрочен
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
