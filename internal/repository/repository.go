package repository

import (
	"context"

	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/models"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *models.Booking) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Booking, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Booking, error)
	Cancel(ctx context.Context, id int64) error
}

