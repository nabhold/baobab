package resolver

import (
	"context"
	"errors"
	"sort"

	"github.com/nabhold/baobab-cp/internal/domain"
)

// MappingResolutionQuery resolves a canonical mapping within a trusted runtime context.
type MappingResolutionQuery struct {
	CanonicalEntityID string
	Context           Context
	Candidates        []domain.Mapping
}

// ResolvedMapping is the chosen mapping together with the specificity score derived from the context scope.
type ResolvedMapping struct {
	Mapping     domain.Mapping
	Specificity int
}

// MappingResolverImpl resolves a mapping candidate list by filtering to active, valid mappings and selecting the highest specificity score.
type MappingResolverImpl struct{}

func (MappingResolverImpl) Resolve(_ context.Context, q MappingResolutionQuery) (ResolvedMapping, error) {
	if q.CanonicalEntityID == "" {
		return ResolvedMapping{}, errors.New("canonical_entity_id is required")
	}
	if len(q.Candidates) == 0 {
		return ResolvedMapping{}, errors.New("mapping not found")
	}

	eligible := make([]domain.Mapping, 0, len(q.Candidates))
	for _, m := range q.Candidates {
		if m.CanonicalEntityID != q.CanonicalEntityID {
			continue
		}
		if m.Status != "ACTIVE" {
			continue
		}
		if err := m.Validate(); err != nil {
			continue
		}
		eligible = append(eligible, m)
	}
	if len(eligible) == 0 {
		return ResolvedMapping{}, errors.New("mapping not found")
	}

	for i := range eligible {
		sel := eligible[i]
		if sel.ScopeID == q.Context.TenantID {
			sel.ResolutionPriority += 20
		}
		if sel.ScopeID == q.Context.LegalEntityID {
			sel.ResolutionPriority += 15
		}
		if sel.ScopeID == q.Context.MarketID {
			sel.ResolutionPriority += 10
		}
		if sel.ScopeID == q.Context.CountryCode {
			sel.ResolutionPriority += 5
		}
		eligible[i] = sel
	}

	sort.Slice(eligible, func(i, j int) bool {
		if eligible[i].ResolutionPriority != eligible[j].ResolutionPriority {
			return eligible[i].ResolutionPriority > eligible[j].ResolutionPriority
		}
		if eligible[i].Confidence != eligible[j].Confidence {
			return confidenceRank(eligible[i].Confidence) > confidenceRank(eligible[j].Confidence)
		}
		return eligible[i].ID > eligible[j].ID
	})

	winner := eligible[0]
	return ResolvedMapping{
		Mapping:     winner,
		Specificity: winner.ResolutionPriority,
	}, nil
}

func confidenceRank(confidence string) int {
	switch confidence {
	case "CONFIRMED":
		return 4
	case "PROBABLE":
		return 3
	case "CANDIDATE":
		return 2
	case "REJECTED":
		return 1
	default:
		return 0
	}
}
