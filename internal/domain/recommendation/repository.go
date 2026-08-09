package recommendation

import "context"

type Repository interface {
	GetByUserID(ctx context.Context, userID string) ([]Recommendation, error)
	Save(ctx context.Context, recommendation Recommendation) error
}
