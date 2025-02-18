package controller

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"time"

	"kinotime/internal/models"
	"kinotime/internal/repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(c *gin.Context, userRepo *repository.UserRepository) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	storedPassword, exists := userRepo.AuthenticateUser(c, user.Username, user.Password)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	userId, err := userRepo.GetUserIdByName(c, user.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"user_id":  userId,
		"exp":      time.Now().Add(time.Hour * time.Duration(JwtExp)).Unix(),
	})

	fmt.Println(user.ID)

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie(
		"auth_token",
		tokenString,
		int(time.Hour*1),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}

func HandleLoginFront(c *gin.Context, templates *template.Template) {
	c.Header("Content-Type", "text/html")

	if templates == nil {
		c.String(http.StatusInternalServerError, "Templates not initialized")
		return
	}

	isAuthenticated, _ := c.Get("isAuthenticated")
	username, _ := c.Get("username")

	templateData := gin.H{
		"IsAuthenticated": isAuthenticated,
	}

	if isAuth, ok := isAuthenticated.(bool); ok && isAuth {
		if usernameStr, ok := username.(string); ok {
			templateData["Username"] = usernameStr
		}
	}

	err := templates.ExecuteTemplate(c.Writer, "login.html", templateData)

	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func HandleRegister(c *gin.Context, userRepo *repository.UserRepository) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		slog.Error("Invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	slog.Info("Attempting to register user", "username", user.Username)

	if _, exists := userRepo.GetUserByUsername(c, user.Username); exists {
		slog.Info("Username already exists", "username", user.Username)
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	err := userRepo.CreateUser(c, user.Username, user.Password)
	if err != nil {
		slog.Error("Error registering user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	slog.Info("User registered successfully", "username", user.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully registered",
	})
}

func HandleRegisterFront(c *gin.Context, templates *template.Template) {
	c.Header("Content-Type", "text/html")

	if templates == nil {
		c.String(http.StatusInternalServerError, "Templates not initialized")
		return
	}

	err := templates.ExecuteTemplate(c.Writer, "register.html", gin.H{})

	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func HandleProfile(c *gin.Context, userRepo *repository.UserRepository) {
	username, _ := c.Get("username")
	password, _ := userRepo.GetUserByUsername(c, username.(string))
	c.JSON(http.StatusOK, gin.H{"username": username, "password": password})
}
