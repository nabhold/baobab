package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nabhold/baobab-cp/internal/domain"
	"github.com/nabhold/baobab-cp/internal/resolver"
)

// PostgresRepository is the PostgreSQL-backed repository implementation for mapping, capability, and topology data.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

var _ MappingRepository = (*PostgresRepository)(nil)
var _ CapabilityRepository = (*PostgresRepository)(nil)
var _ ResolverRepository = (*PostgresRepository)(nil)

func Open(ctx context.Context, url string) (*PostgresRepository, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &PostgresRepository{pool: pool}, nil
}

func (r *PostgresRepository) Close() {
	if r != nil && r.pool != nil {
		r.pool.Close()
	}
}

func (r *PostgresRepository) ListMappings(ctx context.Context, canonicalEntityID string) ([]domain.Mapping, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("repository is not initialized")
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			id,
			mapping_type,
			resolution_mode,
			canonical_entity_id,
			external_reference_id,
			target_canonical_entity_id,
			scope_id,
			direction,
			cardinality,
			authority,
			confidence,
			resolution_priority,
			status,
			effective_from,
			effective_to,
			metadata,
			mapping_version
		FROM mapping.mapping
		WHERE canonical_entity_id = $1
		ORDER BY resolution_priority DESC, confidence DESC`, canonicalEntityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Mapping
	for rows.Next() {
		var m domain.Mapping
		var metadata []byte
		var effectiveFrom string
		var effectiveTo *string
		var externalRef, targetEntity *string
		if err := rows.Scan(
			&m.ID,
			&m.MappingType,
			&m.ResolutionMode,
			&m.CanonicalEntityID,
			&externalRef,
			&targetEntity,
			&m.ScopeID,
			&m.Direction,
			&m.Cardinality,
			&m.Authority,
			&m.Confidence,
			&m.ResolutionPriority,
			&m.Status,
			&effectiveFrom,
			&effectiveTo,
			&metadata,
			&m.Version,
		); err != nil {
			return nil, err
		}
		if externalRef != nil {
			m.ExternalReferenceID = *externalRef
		}
		if targetEntity != nil {
			m.TargetCanonicalEntityID = *targetEntity
		}
		if effectiveTo != nil {
			m.EffectiveTo = *effectiveTo
		}
		m.EffectiveFrom = effectiveFrom
		if len(metadata) > 0 {
			if err := json.Unmarshal(metadata, &m.Metadata); err != nil {
				return nil, fmt.Errorf("decode mapping metadata: %w", err)
			}
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *PostgresRepository) ListBindings(ctx context.Context, capabilityKey string) ([]resolver.CapabilityBinding, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("repository is not initialized")
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			cb.id,
			cap.code,
			ei.engine_id,
			ei.engine_instance_id,
			cb.binding_mode,
			cb.priority,
			cb.status,
			cb.contract_version,
			cb.configuration
		FROM capability.capability_binding cb
		JOIN capability.capability cap ON cap.capability_id = cb.capability_id
		JOIN topology.engine_instance ei ON ei.engine_instance_id = cb.engine_instance_id
		WHERE cap.code = $1
		ORDER BY cb.priority DESC, cb.binding_mode ASC`, capabilityKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []resolver.CapabilityBinding
	for rows.Next() {
		var b resolver.CapabilityBinding
		var config []byte
		if err := rows.Scan(
			&b.CapabilityKey,
			&b.CapabilityKey,
			&b.EngineID,
			&b.EngineInstanceID,
			&b.BindingMode,
			&b.Priority,
			&b.Status,
			&b.ContractVersion,
			&config,
		); err != nil {
			return nil, err
		}
		_ = config
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *PostgresRepository) ListActiveInstances(ctx context.Context, engineID string) ([]resolver.EngineInstance, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("repository is not initialized")
	}
	rows, err := r.pool.Query(ctx, `
		SELECT engine_instance_id, engine_id, region, environment, status
		FROM topology.engine_instance
		WHERE engine_id = $1 AND status = 'ACTIVE'
		ORDER BY region ASC`, engineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []resolver.EngineInstance
	for rows.Next() {
		var instance resolver.EngineInstance
		if err := rows.Scan(&instance.ID, &instance.EngineID, &instance.Region, &instance.Environment, &instance.Status); err != nil {
			return nil, err
		}
		out = append(out, instance)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *PostgresRepository) Ping(ctx context.Context) error {
	if r == nil || r.pool == nil {
		return errors.New("repository is not initialized")
	}
	return r.pool.Ping(ctx)
}

var _ = pgx.ErrNoRows
