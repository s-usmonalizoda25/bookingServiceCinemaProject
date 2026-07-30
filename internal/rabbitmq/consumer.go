package rabbitmq

import "go.uber.org/zap"

func (c *Client) StartConsumer() error {
	msgs, err := c.Channel.Consume(
		"booking_queue",
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			c.Log.Info("received booking event", zap.String("body", string(msg.Body)))
		}
	}()

	return nil
}