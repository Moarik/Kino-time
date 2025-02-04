package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	"kinotime/internal/model"
	"kinotime/internal/types"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	BookingRepo *model.BookingRepository
}

func NewBookingHandler(repo *model.BookingRepository) *BookingHandler {
	return &BookingHandler{BookingRepo: repo}
}

func (h *BookingHandler) HandleCreateBooking(c *gin.Context) {
	var booking types.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		slog.Error(err.Error())
		return
	}

	err := h.BookingRepo.CreateBooking(c, booking.UserID, booking.MovieID, booking.SeatsBooked, booking.TotalPrice, booking.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully"})
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

func (h *BookingHandler) HandleGetAllBookings(c *gin.Context) {
	bookings, err := h.BookingRepo.GetAllBookings(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		slog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

func (h *BookingHandler) HandleUpdateBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var booking types.Booking
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
