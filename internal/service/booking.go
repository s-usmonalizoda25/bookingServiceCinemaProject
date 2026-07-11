package service

import (
	"context"
	"fmt"

	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/gateway/movie"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/gateway/user"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/models"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/repository"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/pkg/errs"
	"go.uber.org/zap"
)

type BookingService struct {
	repo         repository.BookingRepository
	userGateway  user.UserGateway
	movieGateway movie.MovieGateway
	log          *zap.Logger
}

func NewBookingService(repo repository.BookingRepository, userGw user.UserGateway, movieGw movie.MovieGateway, log *zap.Logger) *BookingService {
	return &BookingService{
		repo:         repo,
		userGateway:  userGw,
		movieGateway: movieGw,
		log:          log,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID, movieID int64) (*models.Booking, error) {
	_, err := s.userGateway.GetUser(ctx, userID)
	if err != nil {
		s.log.Error("user validation failed", zap.Int64("userID", userID), zap.Error(err))
		return nil, fmt.Errorf("user not found: %w", err)
	}

	_, err = s.movieGateway.GetMovie(ctx, movieID)
	if err != nil {
		s.log.Error("movie validation failed", zap.Int64("movieID", movieID), zap.Error(err))
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	booking := &models.Booking{
		UserID:  userID,
		MovieID: movieID,
		Status:  models.StatusPending,
	}

	id, err := s.repo.Create(ctx, booking)
	if err != nil {
		s.log.Error(errs.MsgFailedCreate, zap.Error(err))
		return nil, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	booking.ID = id
	return booking, nil
}

func (s *BookingService) GetBooking(ctx context.Context, id int64) (*models.Booking, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error(errs.MsgFailedGet, zap.Int64("id", id), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	return m, nil
}

func (s *BookingService) GetUserBookings(ctx context.Context, userID int64) ([]*models.Booking, error) {
	bookings, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.log.Error(errs.MsgFailedGetUserBooking, zap.Int64("userID", userID), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	return bookings, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, id int64) error {
	err := s.repo.Cancel(ctx, id)
	if err != nil {
		s.log.Error(errs.MsgFailedCancel, zap.Int64("id", id), zap.Error(err))
		return fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	return nil
}
