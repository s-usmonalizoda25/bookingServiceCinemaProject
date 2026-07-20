package models

import "time"

type BookingStatus int32

const (
	StatusUnspecified BookingStatus = 0
	StatusPending     BookingStatus = 1
	StatusConfirmed   BookingStatus = 2
	StatusCancelled   BookingStatus = 3
)

type Booking struct {
	ID        int64
	UserID    int64
	MovieID   int64
	Status    BookingStatus
	CreatedAt time.Time
	UpdatedAt  time.Time
}
