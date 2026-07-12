package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/db"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/gateway/movie"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/gateway/user"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/logger"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/repository"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/server"
	"github.com/s-usmonalizoda25/bookingServiceCinemaProject/internal/service"
	pb "github.com/s-usmonalizoda25/protoCinemaService/gen/booking"
)

func main() {
	if err := godotenv.Load("config/config.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	myLogger := logger.New()
	defer myLogger.Sync()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
	)

	dbPool, err := db.New(context.Background(), dsn)
	if err != nil {
		myLogger.Fatal("failed to connect to db", zap.Error(err))
	}
	defer dbPool.Close()

	userGw, err := user.New(os.Getenv("USER_SERVICE_ADDR"))
	if err != nil {
		myLogger.Fatal("failed to init user gateway", zap.Error(err))
	}

	movieGw, err := movie.New(os.Getenv("MOVIE_SERVICE_ADDR"))
	if err != nil {
		myLogger.Fatal("failed to init movie gateway", zap.Error(err))
	}

	repo := repository.NewBookingRepository(dbPool)
	svc := service.NewBookingService(repo, userGw, movieGw, myLogger.Logger)
	bookingServer := server.New(myLogger.Logger, svc)

	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, bookingServer)
	reflection.Register(grpcServer)

	port := os.Getenv("SERVER_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		myLogger.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		myLogger.Info("booking server started", zap.String("port", port))
		if err := grpcServer.Serve(lis); err != nil {
			myLogger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	myLogger.Info("shutting down server...")
	grpcServer.GracefulStop()
	myLogger.Info("server stopped")
}
