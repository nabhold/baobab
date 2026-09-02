package resolver

import (
	"context"
	"testing"

	"github.com/nabhold/baobab-cp/internal/domain"
)

func TestMappingResolverSelectsMostSpecificCandidate(t *testing.T) {
	resolver := MappingResolverImpl{}
	q := MappingResolutionQuery{
		CanonicalEntityID: "tenant-123",
		Context: Context{
			TenantID:      "tenant-123",
			LegalEntityID: "legal-456",
			MarketID:      "market-789",
			CountryCode:   "ZA",
			CurrencyCode:  "ZAR",
			Locale:        "en-ZA",
		},
		Candidates: []domain.Mapping{
			{
				ID:                      "mapping-market",
				MappingType:             "IDENTITY",
				ResolutionMode:          "SINGLE",
				CanonicalEntityID:       "tenant-123",
				TargetCanonicalEntityID: "entity-market",
				ScopeID:                 "market-789",
				Direction:               "BIDIRECTIONAL",
				Cardinality:             "ONE_TO_ONE",
				Authority:               "baobab",
				Confidence:              "CONFIRMED",
				Status:                  "ACTIVE",
				ResolutionPriority:      20,
				EffectiveFrom:           "2025-01-01T00:00:00Z",
			},
			{
				ID:                      "mapping-tenant",
				MappingType:             "IDENTITY",
				ResolutionMode:          "SINGLE",
				CanonicalEntityID:       "tenant-123",
				TargetCanonicalEntityID: "entity-tenant",
				ScopeID:                 "tenant-123",
				Direction:               "BIDIRECTIONAL",
				Cardinality:             "ONE_TO_ONE",
				Authority:               "baobab",
				Confidence:              "CONFIRMED",
				Status:                  "ACTIVE",
				ResolutionPriority:      50,
				EffectiveFrom:           "2025-01-01T00:00:00Z",
			},
		},
	}

	resolved, err := resolver.Resolve(context.Background(), q)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if resolved.Mapping.ID != "mapping-tenant" {
		t.Fatalf("expected tenant-scoped mapping, got %q", resolved.Mapping.ID)
	}
	if resolved.Specificity <= 0 {
		t.Fatal("expected positive specificity")
	}
}

func TestMappingResolverRejectsInactiveMappings(t *testing.T) {
	resolver := MappingResolverImpl{}
	_, err := resolver.Resolve(context.Background(), MappingResolutionQuery{
		CanonicalEntityID: "tenant-123",
		Context: Context{TenantID: "tenant-123"},
		Candidates: []domain.Mapping{{
			ID:                      "mapping-inactive",
			MappingType:             "IDENTITY",
			ResolutionMode:          "SINGLE",
			CanonicalEntityID:       "tenant-123",
			TargetCanonicalEntityID: "entity-inactive",
			ScopeID:                 "tenant-123",
			Direction:               "BIDIRECTIONAL",
			Cardinality:             "ONE_TO_ONE",
			Authority:               "baobab",
			Confidence:              "CONFIRMED",
			Status:                  "INACTIVE",
			EffectiveFrom:           "2025-01-01T00:00:00Z",
		}},
	})
	if err == nil {
		t.Fatal("expected inactive mapping rejection")
	}
}
