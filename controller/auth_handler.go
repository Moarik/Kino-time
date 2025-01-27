package controller

import (
	"net/http"
	"time"

	"kinotime/internal/model"
	"kinotime/internal/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	jwtSecret = []byte("21cb87e4067d33009b081f3a8090c16d5ae10171a6ca4ddf4b494ad55164a429")
	jwtExp    = 1440
)

func HandleLogin(c *gin.Context, userRepo *model.UserRepository) {
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	storedPassword, exists := userRepo.AuthenticateUser(c, user.Username, user.Password)
	if !exists || storedPassword != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(jwtExp)).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func HandleRegister(c *gin.Context, userRepo *model.UserRepository) {
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if _, exists := userRepo.GetUserByUsername(c, user.Username); exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	err := userRepo.CreateUser(c, user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully registered",
	})
}
