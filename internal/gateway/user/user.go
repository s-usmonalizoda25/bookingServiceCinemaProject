package user

import (
	"context"
	"fmt"

	userpb "github.com/s-usmonalizoda25/protoCinemaService/gen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type UserGateway interface {
	GetUser(ctx context.Context, id int64) (*userpb.GetUserResponse, error)
}

type gateway struct {
	client userpb.UserServiceClient
}

func New(address string) (UserGateway, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	return &gateway{
		client: userpb.NewUserServiceClient(conn),
	}, nil
}

func forwardAuth(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func (g *gateway) GetUser(ctx context.Context, id int64) (*userpb.GetUserResponse, error) {
	return g.client.GetByID(forwardAuth(ctx), &userpb.GetUserRequest{Id: id})
}
