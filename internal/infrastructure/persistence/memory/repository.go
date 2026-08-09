package memory

import (
	"context"
	"sync"

	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

type MemoryRepository struct {
	mu             sync.RWMutex
	recommendations map[string][]recommendation.Recommendation // key: userID
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		recommendations: make(map[string][]recommendation.Recommendation),
	}
}

func (m *MemoryRepository) GetByUserID(
	ctx context.Context,
	userID string,
) ([]recommendation.Recommendation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	recs, exists := m.recommendations[userID]
	if !exists {
		return []recommendation.Recommendation{}, nil
	}

	result := make([]recommendation.Recommendation, len(recs))
	copy(result, recs)

	return result, nil
}

func (m *MemoryRepository) Save(
	ctx context.Context,
	rec recommendation.Recommendation,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	userRecs, exists := m.recommendations[rec.UserID]
	if !exists {
		m.recommendations[rec.UserID] = []recommendation.Recommendation{rec}
		return nil
	}

	for i, existing := range userRecs {
		if existing.ProductID == rec.ProductID {
			userRecs[i] = rec
			return nil
		}
	}

	m.recommendations[rec.UserID] = append(userRecs, rec)
	return nil
}
