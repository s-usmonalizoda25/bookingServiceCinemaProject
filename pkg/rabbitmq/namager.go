package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const prefetchCount = 10

type Manager struct {
	url string

	consumers []Consumer
}

func NewManager(url string) *Manager {
	return &Manager{
		url: url,
	}
}

func (m *Manager) Register(queue string, handler Handler) {
	m.consumers = append(m.consumers, Consumer{
		Queue:   queue,
		Handler: handler,
	})
}

func (m *Manager) Start(ctx context.Context) {
	for {
		if err := m.run(ctx); err != nil {
			log.Println(err)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			log.Println("reconnecting rabbitmq")
		}
	}
}

func (m *Manager) run(ctx context.Context) error {
	conn, err := amqp.Dial(m.url)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}
	defer conn.Close()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for _, consumer := range m.consumers {
		wg.Add(1)
		go func(c Consumer) {
			defer wg.Done()

			for {
				if err := m.consume(runCtx, conn, c); err != nil {
					log.Printf("consumer %s: stopped: %v", c.Queue, err)
				}

				select {
				case <-runCtx.Done():
					return
				case <-time.After(time.Second):
					log.Printf("consumer %s: reconnecting...", c.Queue)
				}
			}
		}(consumer)
	}

	closeErr := <-conn.NotifyClose(make(chan *amqp.Error))

	cancel()
	wg.Wait()

	if closeErr != nil {
		return fmt.Errorf("connection closed: %w", closeErr)
	}
	return nil
}

func (m *Manager) consume(ctx context.Context, conn *amqp.Connection, consumer Consumer) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel: %w", err)
	}
	defer ch.Close()

	if err := ch.Qos(prefetchCount, 0, false); err != nil {
		return fmt.Errorf("ch.Qos: %w", err)
	}

	_, err = ch.QueueDeclare(
		consumer.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("ch.QueueDeclare: %w", err)
	}

	msgs, err := ch.Consume(
		consumer.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("ch.Consume: %w", err)
	}

	notifyClose := ch.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-notifyClose:
			return fmt.Errorf("%s channel closed: %w", consumer.Queue, err)

		case msg, ok := <-msgs:
			if !ok {
				return errors.New("msgs channel closed")
			}

			if err := consumer.Handler(ctx, msg.Body); err != nil {
				if nackErr := msg.Nack(false, false); nackErr != nil {
					log.Printf("consumer %s: nack: %v", consumer.Queue, nackErr)
				}
				continue
			}

			if ackErr := msg.Ack(false); ackErr != nil {
				log.Printf("consumer %s: ack: %v", consumer.Queue, ackErr)
			}
		}
	}
}
