package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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
func (s *Store) RegisterTenant(ctx context.Context, key string, c domain.RegisterTenant) (domain.Operation, error) {
	payload, _ := json.Marshal(c)
	sum := sha256.Sum256(payload)
	hash := hex.EncodeToString(sum[:])
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.Operation{}, err
	}
	defer tx.Rollback(ctx)
	var op domain.Operation
	var prior string
	err = tx.QueryRow(ctx, `SELECT operation_id::text,tenant_id,state,revision,request_hash FROM provisioning_operations WHERE idempotency_key=$1`, key).Scan(&op.OperationID, &op.TenantID, &op.State, &op.Revision, &prior)
	if err == nil {
		if prior != hash {
			return domain.Operation{}, basestore.ErrIdempotencyConflict
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
	if err = tx.QueryRow(ctx, `INSERT INTO provisioning_operations(tenant_id,idempotency_key,request_hash)VALUES($1,$2,$3)RETURNING operation_id::text,tenant_id,state,revision`, c.TenantID, key, hash).Scan(&op.OperationID, &op.TenantID, &op.State, &op.Revision); err != nil {
		return domain.Operation{}, err
	}
	event, _ := json.Marshal(map[string]any{"operation_id": op.OperationID, "tenant_id": c.TenantID, "revision": op.Revision})
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(aggregate_id,event_type,payload)VALUES($1,'tenant.provisioning.started',$2)`, c.TenantID, event); err != nil {
		return domain.Operation{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(tenant_id,action,payload)VALUES($1,'tenant.registration.accepted',$2)`, c.TenantID, payload); err != nil {
		return domain.Operation{}, err
	}
	return op, tx.Commit(ctx)
}
