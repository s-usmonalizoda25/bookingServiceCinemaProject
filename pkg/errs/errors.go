package errs

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
)

const (
	MsgFailedCreate = "failed to create booking"
	MsgFailedGet    = "failed to get booking"
	MsgFailedGetUserBooking = "failed to get user booking"
	MsgFailedCancel = "failed to cancel booking"
)




