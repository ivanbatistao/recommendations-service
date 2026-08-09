package recommendation

import "context"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByUserID(
	ctx context.Context,
	userID string,
) ([]Recommendation, error) {
	return s.repository.GetByUserID(ctx, userID)
}
