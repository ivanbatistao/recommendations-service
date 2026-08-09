package dto

import (
	"testing"

	"github.com/ivanbatistao/recommendations-service/internal/domain/recommendation"
)

func TestFromDomain(t *testing.T) {
	domainRec := recommendation.Recommendation{
		UserID:    "user-123",
		ProductID: "product-456",
		Score:     0.95,
	}

	dto := FromDomain(domainRec)

	if dto.UserID != domainRec.UserID {
		t.Fatalf("expected UserID %s, got %s", domainRec.UserID, dto.UserID)
	}

	if dto.ProductID != domainRec.ProductID {
		t.Fatalf("expected ProductID %s, got %s", domainRec.ProductID, dto.ProductID)
	}

	if dto.Score != domainRec.Score {
		t.Fatalf("expected Score %.2f, got %.2f", domainRec.Score, dto.Score)
	}
}

func TestToDomain(t *testing.T) {
	dto := RecommendationDTO{
		UserID:    "user-123",
		ProductID: "product-456",
		Score:     0.95,
	}

	domainRec := ToDomain(dto)

	if domainRec.UserID != dto.UserID {
		t.Fatalf("expected UserID %s, got %s", dto.UserID, domainRec.UserID)
	}

	if domainRec.ProductID != dto.ProductID {
		t.Fatalf("expected ProductID %s, got %s", dto.ProductID, domainRec.ProductID)
	}

	if domainRec.Score != dto.Score {
		t.Fatalf("expected Score %.2f, got %.2f", dto.Score, domainRec.Score)
	}
}

func TestFromDomainSlice(t *testing.T) {
	domainRecs := []recommendation.Recommendation{
		{
			UserID:    "user-123",
			ProductID: "product-456",
			Score:     0.95,
		},
		{
			UserID:    "user-123",
			ProductID: "product-789",
			Score:     0.82,
		},
	}

	dtos := FromDomainSlice(domainRecs)

	if len(dtos) != len(domainRecs) {
		t.Fatalf("expected %d dtos, got %d", len(domainRecs), len(dtos))
	}

	if dtos[0].ProductID != "product-456" {
		t.Fatalf("expected first product-456, got %s", dtos[0].ProductID)
	}

	if dtos[1].ProductID != "product-789" {
		t.Fatalf("expected second product-789, got %s", dtos[1].ProductID)
	}
}
