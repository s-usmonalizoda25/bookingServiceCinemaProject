package rabbitmq

import "context"

type Handler func(ctx context.Context, body []byte) error

type Consumer struct {
	Queue   string
	Handler Handler
}
