package store

import (
	"context"
	"errors"

	"github.com/nabhold/baobab-cp/internal/domain"
)

var ErrIdempotencyConflict = errors.New("idempotency key was already used with a different request")
var ErrContextDenied = errors.New("tenant context could not be resolved")

type TenantStore interface {
	RegisterTenant(context.Context, string, RequestMetadata, domain.RegisterTenant) (domain.Operation, error)
	ResolveContext(context.Context, RequestMetadata, string, string) (domain.ResolvedContext, error)
	Ping(context.Context) error
}

type RequestMetadata struct {
	ActorID       string
	ActorType     string
	ClientID      string
	TokenID       string
	CorrelationID string
}
