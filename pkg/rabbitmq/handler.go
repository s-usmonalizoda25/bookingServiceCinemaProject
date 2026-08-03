package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type BookingConfirmer interface {
	ConfirmBooking(ctx context.Context, id int64) error
}

type bookingEvent struct {
	Booking struct {
		ID      int64 `json:"id"`
		UserID  int64 `json:"user_id"`
		MovieID int64 `json:"movie_id"`
		Status  int32 `json:"status"`
	} `json:"booking"`
}

func RegisterConsumers(manager *Manager, confirmer BookingConfirmer, log *zap.Logger) {
	manager.Register("booking_queue", func(ctx context.Context, body []byte) error {
		var event bookingEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Error("failed to unmarshal booking event", zap.Error(err), zap.String("body", string(body)))
			return fmt.Errorf("unmarshal booking event: %w", err)
		}

		log.Info("received booking event in consumer",
			zap.Int64("booking_id", event.Booking.ID),
			zap.Int64("user_id", event.Booking.UserID),
		)

		if err := confirmer.ConfirmBooking(ctx, event.Booking.ID); err != nil {
			return fmt.Errorf("confirm booking %d: %w", event.Booking.ID, err)
		}

		return nil
	})
}
