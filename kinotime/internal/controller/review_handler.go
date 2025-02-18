package controller

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"kinotime/internal/models"
	"kinotime/internal/repository"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	ReviewRepo *repository.ReviewRepository
	Templates  *template.Template
	MovieRepo  *repository.MovieRepository
}

func NewReviewHandler(reviewrepo *repository.ReviewRepository, movieRepo *repository.MovieRepository, templates *template.Template) *ReviewHandler {
	return &ReviewHandler{ReviewRepo: reviewrepo, MovieRepo: movieRepo, Templates: templates}
}

func (h *ReviewHandler) HandleCreateReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		slog.Error(err.Error())
		return
	}

	err := h.ReviewRepo.CreateReview(c, review.UserID, review.MovieID, review.Rating, review.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review created successfully"})
}

func (h *ReviewHandler) HandleGetReviewFront(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username := c.Param("username")

	movieIDStr := c.Param("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	reviews, err := h.ReviewRepo.GetReviewsByMovieID(c, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		slog.Error(err.Error())
		return
	}

	movie, err := h.MovieRepo.GetMovieByID(c, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie"})
		slog.Error(err.Error())
		return
	}

	err = h.Templates.ExecuteTemplate(c.Writer, "review.html", gin.H{
		"UserID":          userIDStr,
		"IsAuthenticated": true,
		"Username":        username,
		"MovieID":         movieID,
		"MovieTitle":      movie.Title,
		"Reviews":         reviews,
	})

	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func (h *ReviewHandler) HandleCreateReviewForm(c *gin.Context) {
	// Get user ID from session - note that it's stored as a string
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert the userID from string to int
	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse user ID"})
		return
	}

	// Parse form data
	movieIDStr := c.PostForm("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	ratingStr := c.PostForm("rating")
	rating, err := strconv.Atoi(ratingStr)
	if err != nil || rating < 1 || rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating"})
		return
	}

	comment := c.PostForm("comment")

	// Create the review
	err = h.ReviewRepo.CreateReview(c, userID, movieID, rating, comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		slog.Error(err.Error())
		return
	}

	// Redirect back to the reviews page
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/front/review/%d", movieID))
}

func (h *ReviewHandler) HandleGetReviewByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.ReviewRepo.GetReviewByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"review": review})
}

func (h *ReviewHandler) HandleGetReviewsByMovieID(c *gin.Context) {
	movieIDStr := c.Param("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	reviews, err := h.ReviewRepo.GetReviewsByMovieID(c, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (h *ReviewHandler) HandleUpdateReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		slog.Error(err.Error())
		return
	}

	err = h.ReviewRepo.UpdateReview(c, id, review.Rating, review.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
}

func (h *ReviewHandler) HandleDeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	err = h.ReviewRepo.DeleteReview(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
