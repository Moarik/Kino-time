package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kinotime/internal/models"
)

type BookingRepository struct {
	DB *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{DB: db}
}

func (repo *BookingRepository) CreateBooking(ctx context.Context, userID, movieID, seatsBooked int, totalPrice float64, status, bookingTime string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := repo.DB.ExecContext(ctx, `
		INSERT INTO bookings (user_id, movie_id, seats_booked, total_price, status, booking_time, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID, movieID, seatsBooked, totalPrice, status, bookingTime, now, now)
	return err
}

func (repo *BookingRepository) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	var booking models.Booking
	err := repo.DB.QueryRowContext(ctx, "SELECT id, user_id, movie_id, seats_booked, total_price, status FROM bookings WHERE id = $1", id).
		Scan(&booking.ID, &booking.UserID, &booking.MovieID, &booking.SeatsBooked, &booking.TotalPrice, &booking.Status)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (repo *BookingRepository) GetBookingsByUserID(ctx context.Context, userID string) ([]models.Booking, error) {
	query := `
        SELECT b.id, b.user_id, b.movie_id, m.title as movie_title,
               b.seats_booked, b.total_price, b.status, b.booking_time,
               b.created_at, b.updated_at
        FROM bookings b
        JOIN movies m ON b.movie_id = m.id
        WHERE b.user_id = $1
        ORDER BY b.created_at DESC`

	rows, err := repo.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying bookings: %v", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.MovieID,
			&booking.MovieTitle,
			&booking.SeatsBooked,
			&booking.TotalPrice,
			&booking.Status,
			&booking.BookingTime,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning booking: %v", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return bookings, nil
}

func (repo *BookingRepository) UpdateBooking(ctx context.Context, id int, seatsBooked int, totalPrice float64, status string) error {
	_, err := repo.DB.ExecContext(ctx, "UPDATE bookings SET seats_booked = $1, total_price = $2, status = $3 WHERE id = $4", seatsBooked, totalPrice, status, id)
	return err
}

func (repo *BookingRepository) DeleteBooking(ctx context.Context, id int) error {
	_, err := repo.DB.ExecContext(ctx, "DELETE FROM bookings WHERE id = $1", id)
	return err
}
