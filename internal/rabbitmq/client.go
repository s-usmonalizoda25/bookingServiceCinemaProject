package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Client struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Log     *zap.Logger
}

func New(url string, log *zap.Logger) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a rabbitmq channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		"booking_queue", // имя очереди
		true,            // durable
		false,           // auto-deleted
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	log.Info("successfully connected to rabbitmq and declared queue")

	return &Client{
		Conn:    conn,
		Channel: ch,
		Log:     log,
	}, nil
}

func (c *Client) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}
