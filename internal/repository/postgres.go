package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/models"
)

type bookingRepo struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) BookingRepository {
	return &bookingRepo{db: db}
}

func (r *bookingRepo) Create(ctx context.Context, b *models.Booking) (int64, error) {
	query := `INSERT INTO bookings (user_id, movie_id, status) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.db.QueryRow(ctx, query, b.UserID, b.MovieID, b.Status).Scan(&id)
	return id, err
}

func (r *bookingRepo) GetByID(ctx context.Context, id int64) (*models.Booking, error) {
	query := `SELECT id, user_id, movie_id, status, created_at, updated_at FROM bookings WHERE id = $1`
	b := &models.Booking{}
	err := r.db.QueryRow(ctx, query, id).Scan(&b.ID, &b.UserID, &b.MovieID, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

func (r *bookingRepo) GetByUserID(ctx context.Context, userID int64) ([]*models.Booking, error) {
	query := `SELECT id, user_id, movie_id, status, created_at, updated_at FROM bookings WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		b := &models.Booking{}
		if err := rows.Scan(&b.ID, &b.UserID, &b.MovieID, &b.Status, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *bookingRepo) Cancel(ctx context.Context, id int64) error {
	query := `UPDATE bookings SET status = 3, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
