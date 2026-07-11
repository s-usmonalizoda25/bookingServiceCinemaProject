package movie

import (
	"context"
	"fmt"

	moviepb "github.com/s-usmonalizoda25/movieServiceCinemaProject/genproto/movie"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MovieGateway interface {
	GetMovie(ctx context.Context, id int64) (*moviepb.GetMovieResponse, error)
}

type gateway struct {
	client moviepb.MovieServiceClient
}

func New(address string) (MovieGateway, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to movie service: %w", err)
	}

	return &gateway{
		client: moviepb.NewMovieServiceClient(conn),
	}, nil
}

func (g *gateway) GetMovie(ctx context.Context, id int64) (*moviepb.GetMovieResponse, error) {
	return g.client.GetByID(ctx, &moviepb.GetMovieRequest{Id: id})
}
