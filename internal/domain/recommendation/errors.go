package recommendation

import "errors"

var (
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidProductID   = errors.New("invalid product ID")
	ErrInvalidEventType   = errors.New("invalid event type")
	ErrInvalidLimit       = errors.New("invalid limit")
	ErrUserNotFound       = errors.New("user not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrRepositoryError    = errors.New("repository error")
	ErrRecommendationNotFound = errors.New("recommendation not found")
)
