package store

import (
	"context"
	"errors"
	"github.com/nabhold/baobab-cp/internal/domain"
)

var ErrIdempotencyConflict = errors.New("idempotency key was already used with a different request")

type TenantStore interface {
	RegisterTenant(context.Context, string, domain.RegisterTenant) (domain.Operation, error)
	Ping(context.Context) error
}
