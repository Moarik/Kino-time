package controller

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"kinotime/internal/models"
	"kinotime/internal/repository"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	BookingRepo *repository.BookingRepository
	Templates   *template.Template
}

func NewBookingHandler(repo *repository.BookingRepository, templates *template.Template) *BookingHandler {
	return &BookingHandler{BookingRepo: repo, Templates: templates}
}

func (h *BookingHandler) HandleCreateBooking(c *gin.Context) {
	var bookingReq models.BookingRequest

	isAuthenticated, _ := c.Get("isAuthenticated")
	if !isAuthenticated.(bool) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		slog.Error(err.Error())
		return
	}

	movieID, err := strconv.Atoi(c.Request.FormValue("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	seatsBooked, err := strconv.Atoi(c.Request.FormValue("seats_booked"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seats number"})
		return
	}

	totalPrice, err := strconv.ParseFloat(c.Request.FormValue("total_price"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}

	bookingReq = models.BookingRequest{
		MovieID:     movieID,
		SeatsBooked: seatsBooked,
		TotalPrice:  totalPrice,
		Status:      c.Request.FormValue("status"),
		BookingTime: c.Request.FormValue("booking_time"),
	}

	err = h.BookingRepo.CreateBooking(c, userID, bookingReq.MovieID, bookingReq.SeatsBooked,
		bookingReq.TotalPrice, bookingReq.Status, bookingReq.BookingTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		slog.Error(err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (h *BookingHandler) HandleGetBookingByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := h.BookingRepo.GetBookingByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"booking": booking})
}

//
//func (h *BookingHandler) HandleGetAllBookings(c *gin.Context) {
//	bookings, err := h.BookingRepo.GetAllBookings(c)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
//		slog.Error(err.Error())
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
//}

func (h *BookingHandler) HandleUpdateBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		slog.Error(err.Error())
		return
	}

	err = h.BookingRepo.UpdateBooking(c, id, booking.SeatsBooked, booking.TotalPrice, booking.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking updated successfully"})
}

func (h *BookingHandler) HandleDeleteBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	err = h.BookingRepo.DeleteBooking(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}

func (h *BookingHandler) HandleGetBookingPage(c *gin.Context) {
	movieIDStr := c.Param("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.Templates.ExecuteTemplate(c.Writer, "booking.html", gin.H{
		"MovieID":    movieID,
		"UserID":     userId,
		"MoviePrice": 10.00,
	})

	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func (h *BookingHandler) HandleGetBookingUserPage(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	bookings, err := h.BookingRepo.GetBookingsByUserID(c, userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		slog.Error(err.Error())
		return
	}

	err = h.Templates.ExecuteTemplate(c.Writer, "orders.html", gin.H{
		"UserID":     userIDStr,
		"MoviePrice": 10.00,
		"Bookings":   bookings,
	})

	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}
