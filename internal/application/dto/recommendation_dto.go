package dto

import "github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"

type RecommendationDTO struct {
	UserID    string  `json:"user_id"`
	ProductID string  `json:"product_id"`
	Score     float64 `json:"score"`
}

func FromDomain(rec recommendation.Recommendation) RecommendationDTO {
	return RecommendationDTO{
		UserID:    rec.UserID,
		ProductID: rec.ProductID,
		Score:     rec.Score,
	}
}

func ToDomain(dto RecommendationDTO) recommendation.Recommendation {
	return recommendation.Recommendation{
		UserID:    dto.UserID,
		ProductID: dto.ProductID,
		Score:     dto.Score,
	}
}

func FromDomainSlice(recs []recommendation.Recommendation) []RecommendationDTO {
	dtos := make([]RecommendationDTO, len(recs))
	for i, rec := range recs {
		dtos[i] = FromDomain(rec)
	}
	return dtos
}
