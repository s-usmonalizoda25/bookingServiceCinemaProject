package server

import (
	"context"

	pb "github.com/s-usmonalizoda25/protoCinemaService/gen/booking"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/service"
	"go.uber.org/zap"
)

type Server struct {
	pb.UnimplementedBookingServiceServer
	log *zap.Logger
	svc *service.BookingService
}

func New(log *zap.Logger, svc *service.BookingService) *Server {
	return &Server{
		log: log,
		svc: svc,
	}
}

func (s *Server) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	booking, err := s.svc.CreateBooking(ctx, req.UserId, req.MovieId)
	if err != nil {
		return nil, err
	}

	return &pb.CreateBookingResponse{
		Booking: &pb.Booking{
			Id:      booking.ID,
			UserId:  booking.UserID,
			MovieId: booking.MovieID,
			Status:  pb.BookingStatus(booking.Status),
		},
	}, nil
}

func (s *Server) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.GetBookingResponse, error) {
	b, err := s.svc.GetBooking(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetBookingResponse{
		Booking: &pb.Booking{
			Id:      b.ID,
			UserId:  b.UserID,
			MovieId: b.MovieID,
			Status:  pb.BookingStatus(b.Status),
		},
	}, nil
}
