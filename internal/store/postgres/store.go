package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nabhold/baobab-cp/internal/domain"
	basestore "github.com/nabhold/baobab-cp/internal/store"
)

type Store struct{ pool *pgxpool.Pool }

func Open(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Store{pool}, nil
}
func (s *Store) Close()                         { s.pool.Close() }
func (s *Store) Ping(ctx context.Context) error { return s.pool.Ping(ctx) }
func (s *Store) RegisterTenant(ctx context.Context, key string, metadata basestore.RequestMetadata, c domain.RegisterTenant) (domain.Operation, error) {
	requestPayload, _ := json.Marshal(c)
	sum := sha256.Sum256(requestPayload)
	hash := hex.EncodeToString(sum[:])
	auditPayload, _ := json.Marshal(map[string]any{
		"legal_entity_id":    c.LegalEntityID,
		"requested_products": c.RequestedProducts,
		"isolation_strategy": c.IsolationStrategy,
		"residency_region":   c.ResidencyRegion,
	})
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.Operation{}, err
	}
	defer tx.Rollback(ctx)
	var op domain.Operation
	var prior string
	err = tx.QueryRow(ctx, `SELECT operation_id::text,tenant_id,state,revision,created_at,updated_at,request_hash FROM provisioning_operations WHERE idempotency_key=$1`, key).Scan(&op.OperationID, &op.TenantID, &op.State, &op.Revision, &op.CreatedAt, &op.UpdatedAt, &prior)
	if err == nil {
		if prior != hash {
			if err = insertAudit(ctx, tx, op.TenantID, metadata, key, "tenant.registration.requested", "tenant:"+op.TenantID, "denied", "idempotency_conflict", auditPayload); err != nil {
				return domain.Operation{}, err
			}
			if err = tx.Commit(ctx); err != nil {
				return domain.Operation{}, err
			}
			return domain.Operation{}, basestore.ErrIdempotencyConflict
		}
		if err = insertAudit(ctx, tx, op.TenantID, metadata, key, "tenant.registration.requested", "tenant:"+op.TenantID, "accepted", "idempotent_replay", auditPayload); err != nil {
			return domain.Operation{}, err
		}
		return op, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Operation{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO legal_entities(legal_entity_id)VALUES($1)ON CONFLICT DO NOTHING`, c.LegalEntityID); err != nil {
		return domain.Operation{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO tenants(tenant_id,legal_entity_id,display_name,isolation_strategy,residency_region,metadata)VALUES($1,$2,$3,$4,$5,$6)`, c.TenantID, c.LegalEntityID, c.DisplayName, c.IsolationStrategy, c.ResidencyRegion, c.Metadata); err != nil {
		return domain.Operation{}, err
	}
	for _, product := range c.RequestedProducts {
		if _, err = tx.Exec(ctx, `INSERT INTO product_subscriptions(tenant_id,product_id)VALUES($1,$2)`, c.TenantID, product); err != nil {
			return domain.Operation{}, err
		}
	}
	if err = tx.QueryRow(ctx, `INSERT INTO provisioning_operations(tenant_id,idempotency_key,request_hash)VALUES($1,$2,$3)RETURNING operation_id::text,tenant_id,state,revision,created_at,updated_at`, c.TenantID, key, hash).Scan(&op.OperationID, &op.TenantID, &op.State, &op.Revision, &op.CreatedAt, &op.UpdatedAt); err != nil {
		return domain.Operation{}, err
	}
	event, _ := json.Marshal(map[string]any{"operation_id": op.OperationID, "tenant_id": c.TenantID, "revision": op.Revision, "correlation_id": metadata.CorrelationID})
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(aggregate_id,event_type,payload)VALUES($1,'tenant.provisioning.started',$2)`, c.TenantID, event); err != nil {
		return domain.Operation{}, err
	}
	if err = insertAudit(ctx, tx, c.TenantID, metadata, key, "tenant.registration.requested", "tenant:"+c.TenantID, "accepted", "scope_allowed", auditPayload); err != nil {
		return domain.Operation{}, err
	}
	return op, tx.Commit(ctx)
}

func (s *Store) ResolveContext(ctx context.Context, metadata basestore.RequestMetadata, tenantID, productID string) (domain.ResolvedContext, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.ResolvedContext{}, err
	}
	defer tx.Rollback(ctx)

	var legalEntityID, desiredState, observedState, subscriptionStatus string
	err = tx.QueryRow(ctx, `
		SELECT t.legal_entity_id,t.desired_state,t.observed_state,COALESCE(ps.status,'')
		FROM tenants t
		LEFT JOIN product_subscriptions ps ON ps.tenant_id=t.tenant_id AND ps.product_id=$2
		WHERE t.tenant_id=$1`, tenantID, productID).Scan(&legalEntityID, &desiredState, &observedState, &subscriptionStatus)

	found := true
	if errors.Is(err, pgx.ErrNoRows) {
		found = false
	} else if err != nil {
		return domain.ResolvedContext{}, err
	}
	result, policyDecision := contextPolicyDecision(found, desiredState, observedState, subscriptionStatus)

	auditPayload, marshalErr := json.Marshal(map[string]string{
		"product_id":          productID,
		"desired_state":       desiredState,
		"observed_state":      observedState,
		"subscription_status": subscriptionStatus,
	})
	if marshalErr != nil {
		return domain.ResolvedContext{}, fmt.Errorf("marshal context audit: %w", marshalErr)
	}
	target := "tenant:" + tenantID + "/product:" + productID
	if err = insertAudit(ctx, tx, tenantID, metadata, "", "context.resolve", target, result, policyDecision, auditPayload); err != nil {
		return domain.ResolvedContext{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.ResolvedContext{}, err
	}
	if result == "denied" {
		return domain.ResolvedContext{}, basestore.ErrContextDenied
	}

	return domain.ResolvedContext{
		TenantID:        tenantID,
		EntityID:        legalEntityID,
		LifecycleStatus: "active",
		ProductID:       productID,
		Entitled:        true,
		CacheTTLSeconds: 15,
		ResolvedAt:      time.Now().UTC(),
		CorrelationID:   metadata.CorrelationID,
	}, nil
}

func contextPolicyDecision(found bool, desiredState, observedState, subscriptionStatus string) (string, string) {
	if !found {
		return "denied", "tenant_unknown"
	}
	if desiredState != "active" || observedState != "active" {
		return "denied", "tenant_not_active"
	}
	if subscriptionStatus != "active" {
		return "denied", "product_not_entitled"
	}
	return "allowed", "context_allowed"
}

func insertAudit(ctx context.Context, tx pgx.Tx, tenantID string, metadata basestore.RequestMetadata, idempotencyKey, action, target, result, decision string, payload []byte) error {
	_, err := tx.Exec(ctx, `INSERT INTO audit_events(tenant_id,actor_id,actor_type,client_id,token_id,correlation_id,idempotency_key,action,target,result,policy_decision,payload) VALUES($1,$2,$3,$4,$5,$6,NULLIF($7,''),$8,$9,$10,$11,$12)`, tenantID, metadata.ActorID, metadata.ActorType, metadata.ClientID, metadata.TokenID, metadata.CorrelationID, idempotencyKey, action, target, result, decision, payload)
	return err
}
