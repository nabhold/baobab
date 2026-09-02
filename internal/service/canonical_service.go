package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nabhold/baobab-cp/internal/domain"
	"github.com/nabhold/baobab-cp/internal/repository"
)

// CanonicalEntityService owns canonical entity lifecycle transitions.
type CanonicalEntityService struct {
	Repository repository.CanonicalEntityRepository
	Now        func() time.Time
}

func (s CanonicalEntityService) Create(ctx context.Context, entity domain.CanonicalEntity) (domain.CanonicalEntity, error) {
	if s.Repository == nil {
		return domain.CanonicalEntity{}, errors.New("canonical repository is required")
	}
	if entity.Status == "" {
		entity.Status = "DRAFT"
	}
	if entity.SchemaVersion == 0 {
		entity.SchemaVersion = 1
	}
	if entity.EffectiveFrom.IsZero() {
		if s.Now != nil {
			entity.EffectiveFrom = s.Now().UTC()
		} else {
			entity.EffectiveFrom = time.Now().UTC()
		}
	}
	if err := s.Repository.CreateCanonicalEntity(ctx, entity); err != nil {
		return domain.CanonicalEntity{}, err
	}
	entity.Version = 1
	return entity, nil
}

func (s CanonicalEntityService) Get(ctx context.Context, id string) (domain.CanonicalEntity, error) {
	if s.Repository == nil {
		return domain.CanonicalEntity{}, errors.New("canonical repository is required")
	}
	return s.Repository.GetCanonicalEntity(ctx, id)
}

func (s CanonicalEntityService) Validate(ctx context.Context, id string, expectedVersion int64) (domain.CanonicalEntity, error) {
	return s.transition(ctx, id, expectedVersion, "VALIDATED", "DRAFT")
}

func (s CanonicalEntityService) Activate(ctx context.Context, id string, expectedVersion int64) (domain.CanonicalEntity, error) {
	return s.transition(ctx, id, expectedVersion, "ACTIVE", "VALIDATED")
}

func (s CanonicalEntityService) Suspend(ctx context.Context, id string, expectedVersion int64) (domain.CanonicalEntity, error) {
	return s.transition(ctx, id, expectedVersion, "SUSPENDED", "ACTIVE")
}

func (s CanonicalEntityService) Retire(ctx context.Context, id string, expectedVersion int64) (domain.CanonicalEntity, error) {
	entity, err := s.Get(ctx, id)
	if err != nil {
		return domain.CanonicalEntity{}, err
	}
	if entity.Status != "ACTIVE" && entity.Status != "SUSPENDED" && entity.Status != "DEPRECATED" {
		return domain.CanonicalEntity{}, fmt.Errorf("canonical entity %s cannot transition from %s to RETIRED", id, entity.Status)
	}
	return s.transition(ctx, id, expectedVersion, "RETIRED", entity.Status)
}

func (s CanonicalEntityService) transition(ctx context.Context, id string, expectedVersion int64, next, required string) (domain.CanonicalEntity, error) {
	entity, err := s.Get(ctx, id)
	if err != nil {
		return domain.CanonicalEntity{}, err
	}
	if entity.Status != required {
		return domain.CanonicalEntity{}, fmt.Errorf("canonical entity %s cannot transition from %s to %s", id, entity.Status, next)
	}
	entity.Status = next
	if err := s.Repository.SaveCanonicalEntity(ctx, entity, expectedVersion); err != nil {
		return domain.CanonicalEntity{}, err
	}
	entity.Version = expectedVersion + 1
	return entity, nil
}

func ParseExpectedVersion(value string) (int64, error) {
	if strings.TrimSpace(value) == "" {
		return 0, errors.New("If-Match is required")
	}
	var version int64
	if _, err := fmt.Sscan(strings.Trim(value, "\""), &version); err != nil || version < 1 {
		return 0, errors.New("If-Match must contain a positive entity version")
	}
	return version, nil
}
