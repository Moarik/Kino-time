package models

import "time"

type Booking struct {
	ID          int       `json:"id"`
	UserID      string    `json:"user_id"`
	MovieID     int       `json:"movie_id"`
	MovieTitle  string    `json:"movie_title"`
	SeatsBooked int       `json:"seats_booked"`
	TotalPrice  float64   `json:"total_price"`
	Status      string    `json:"status"`
	BookingTime string    `json:"booking_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
