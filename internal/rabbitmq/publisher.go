package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func (c *Client) Publish(body string) error {
	err := c.Channel.Publish(
		"",              // exchange
		"booking_queue", // routing key (имя очереди)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	c.Log.Info("sent message to rabbitmq", zap.String("body", body))
	return nil
}
