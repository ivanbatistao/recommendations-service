package event

import "errors"

var (
	ErrInvalidEventID     = errors.New("invalid event ID")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidProductID   = errors.New("invalid product ID")
	ErrInvalidEventType   = errors.New("invalid event type")
	ErrInvalidTimestamp   = errors.New("invalid timestamp")
	ErrInvalidMetadata    = errors.New("invalid metadata")
	ErrEventNotFound      = errors.New("event not found")
)
