package rabbitmq

import (
    "context"
    "go.uber.org/zap"
)

func RegisterConsumers(manager *Manager, log *zap.Logger) {
    manager.Register("booking_queue", func(ctx context.Context, body []byte) error {
        log.Info("received booking event in consumer", zap.String("body", string(body)))
        
        return nil
    })
}