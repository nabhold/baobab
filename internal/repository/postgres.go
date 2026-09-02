package repository

import (
	"context"
	"errors"
	"fmt"

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
var _ MappingWriter = (*PostgresRepository)(nil)
var _ CapabilityWriter = (*PostgresRepository)(nil)
var _ CanonicalEntityRepository = (*PostgresRepository)(nil)

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

func (r *PostgresRepository) CreateCanonicalEntity(ctx context.Context, entity domain.CanonicalEntity) error {
	if r == nil || r.pool == nil {
		return errors.New("repository is not initialized")
	}
	if err := entity.Validate(); err != nil {
		return fmt.Errorf("validate canonical entity: %w", err)
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO registry.canonical_entity(canonical_entity_id, tenant_id, legal_entity_id, entity_type, external_key, status)
		VALUES ($1::uuid, NULLIF($2, ''), NULLIF($3, ''), $4, NULLIF($5, ''), LOWER($6))`, entity.ID, entity.OwnerTenantID, entity.OwnerLegalEntityID, entity.EntityType, entity.CanonicalKey, entity.Status)
	return err
}

func (r *PostgresRepository) GetCanonicalEntity(ctx context.Context, id string) (domain.CanonicalEntity, error) {
	if r == nil || r.pool == nil {
		return domain.CanonicalEntity{}, errors.New("repository is not initialized")
	}
	var entity domain.CanonicalEntity
	var status string
	err := r.pool.QueryRow(ctx, `SELECT canonical_entity_id::text, COALESCE(tenant_id,''), COALESCE(legal_entity_id,''), entity_type, COALESCE(external_key,''), UPPER(status), version, created_at, updated_at FROM registry.canonical_entity WHERE canonical_entity_id=$1::uuid`, id).Scan(&entity.ID, &entity.OwnerTenantID, &entity.OwnerLegalEntityID, &entity.EntityType, &entity.CanonicalKey, &status, &entity.Version, &entity.CreatedAt, &entity.UpdatedAt)
	if err != nil {
		return domain.CanonicalEntity{}, fmt.Errorf("get canonical entity %s: %w", id, err)
	}
	entity.Status, entity.SchemaVersion, entity.Authority, entity.Classification = status, 1, "baobab", "INTERNAL"
	entity.DisplayName = entity.CanonicalKey
	entity.EffectiveFrom = entity.CreatedAt
	return entity, nil
}

func (r *PostgresRepository) SaveCanonicalEntity(ctx context.Context, entity domain.CanonicalEntity, expectedVersion int64) error {
	if r == nil || r.pool == nil {
		return errors.New("repository is not initialized")
	}
	result, err := r.pool.Exec(ctx, `UPDATE registry.canonical_entity SET entity_type=$2, external_key=NULLIF($3,''), status=LOWER($4), version=version+1, updated_at=now() WHERE canonical_entity_id=$1::uuid AND version=$5`, entity.ID, entity.EntityType, entity.CanonicalKey, entity.Status, expectedVersion)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("canonical entity %s version conflict or not found", entity.ID)
	}
	return nil
}

func (r *PostgresRepository) ListMappings(ctx context.Context, canonicalEntityID string) ([]domain.Mapping, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("repository is not initialized")
	}
	rows, err := r.pool.Query(ctx, `
		SELECT cm.canonical_mapping_id::text, cm.mapping_type, cm.source_entity_id::text,
		       cm.target_entity_id::text, COALESCE(ms.mapping_scope_id::text, cm.source_entity_id::text),
		       UPPER(cm.status), cm.created_at::text
		FROM mapping.canonical_mapping cm
		JOIN registry.canonical_entity ce ON ce.canonical_entity_id = cm.source_entity_id
		LEFT JOIN mapping.mapping_scope ms ON ms.tenant_id = ce.tenant_id
		WHERE ce.tenant_id = $1 OR cm.source_entity_id::text = $1
		ORDER BY cm.created_at DESC`, canonicalEntityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Mapping
	for rows.Next() {
		var m domain.Mapping
		var effectiveFrom string
		if err := rows.Scan(
			&m.ID,
			&m.MappingType,
			&m.CanonicalEntityID,
			&m.TargetCanonicalEntityID,
			&m.ScopeID,
			&m.Status,
			&effectiveFrom,
		); err != nil {
			return nil, err
		}
		m.ResolutionMode = "SINGLE"
		m.Direction = "SOURCE_TO_TARGET"
		m.Cardinality = "ONE_TO_ONE"
		m.Authority = "baobab"
		m.Confidence = "CONFIRMED"
		m.ResolutionPriority = 0
		m.Version = 1
		m.EffectiveFrom = effectiveFrom
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *PostgresRepository) CreateMapping(ctx context.Context, mapping domain.Mapping) error {
	if r == nil || r.pool == nil {
		return errors.New("repository is not initialized")
	}
	if err := mapping.Validate(); err != nil {
		return fmt.Errorf("validate mapping: %w", err)
	}
	if mapping.ExternalReferenceID != "" {
		return errors.New("external-reference mappings are not supported by the current canonical schema")
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO mapping.canonical_mapping(canonical_mapping_id, source_entity_id, target_entity_id, mapping_type, status)
		VALUES ($1::uuid, $2::uuid, $3::uuid, $4, LOWER($5))`, mapping.ID, mapping.CanonicalEntityID, mapping.TargetCanonicalEntityID, mapping.MappingType, mapping.Status)
	return err
}

func (r *PostgresRepository) GetMapping(ctx context.Context, mappingID string) (domain.Mapping, error) {
	if r == nil || r.pool == nil {
		return domain.Mapping{}, errors.New("repository is not initialized")
	}
	var mapping domain.Mapping
	var effectiveFrom string
	err := r.pool.QueryRow(ctx, `SELECT canonical_mapping_id::text, mapping_type, source_entity_id::text, target_entity_id::text, source_entity_id::text, UPPER(status), created_at::text FROM mapping.canonical_mapping WHERE canonical_mapping_id=$1::uuid`, mappingID).Scan(&mapping.ID, &mapping.MappingType, &mapping.CanonicalEntityID, &mapping.TargetCanonicalEntityID, &mapping.ScopeID, &mapping.Status, &effectiveFrom)
	if err != nil {
		return domain.Mapping{}, fmt.Errorf("get mapping %s: %w", mappingID, err)
	}
	mapping.ResolutionMode, mapping.Direction, mapping.Cardinality = "SINGLE", "SOURCE_TO_TARGET", "ONE_TO_ONE"
	mapping.Authority, mapping.Confidence, mapping.EffectiveFrom, mapping.Version = "baobab", "CONFIRMED", effectiveFrom, 1
	return mapping, nil
}

func (r *PostgresRepository) SaveMapping(ctx context.Context, mapping domain.Mapping, expectedVersion int64) error {
	if r == nil || r.pool == nil {
		return errors.New("repository is not initialized")
	}
	if err := mapping.Validate(); err != nil {
		return fmt.Errorf("validate mapping: %w", err)
	}
	if expectedVersion != 1 {
		return fmt.Errorf("mapping %s version conflict: expected %d, got 1", mapping.ID, expectedVersion)
	}
	result, err := r.pool.Exec(ctx, `UPDATE mapping.canonical_mapping SET mapping_type=$2, status=LOWER($3) WHERE canonical_mapping_id=$1::uuid`, mapping.ID, mapping.MappingType, mapping.Status)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("mapping %s not found", mapping.ID)
	}
	return nil
}

func (r *PostgresRepository) ListBindings(ctx context.Context, capabilityKey string) ([]resolver.CapabilityBinding, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("repository is not initialized")
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			cb.id::text,
			cap.code,
			ei.engine_id,
			ei.engine_instance_id,
			cb.binding_mode,
			cb.priority,
			UPPER(cb.status),
			cb.contract_version,
			cb.scope_id::text
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
		if err := rows.Scan(
			&b.ID,
			&b.CapabilityKey,
			&b.EngineID,
			&b.EngineInstanceID,
			&b.BindingMode,
			&b.Priority,
			&b.Status,
			&b.ContractVersion,
			&b.ScopeID,
		); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *PostgresRepository) CreateBinding(ctx context.Context, binding resolver.CapabilityBinding) error {
	if r == nil || r.pool == nil {
		return errors.New("repository is not initialized")
	}
	if binding.CapabilityKey == "" || binding.EngineID == "" || binding.EngineInstanceID == "" || binding.ScopeID == "" {
		return errors.New("capability, engine, engine instance, and scope are required")
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO capability.capability_binding(capability_id, engine_instance_id, scope_id, binding_mode, priority, status, contract_version, effective_from)
		SELECT c.capability_id, ei.engine_instance_id, $4::uuid, $5, $6, LOWER($7), $8, now()
		FROM capability.capability c JOIN topology.engine_instance ei ON ei.engine_instance_id=$3::uuid AND ei.engine_id=$2::uuid
		WHERE c.code=$1`, binding.CapabilityKey, binding.EngineID, binding.EngineInstanceID, binding.ScopeID, binding.BindingMode, binding.Priority, binding.Status, binding.ContractVersion)
	return err
}

func (r *PostgresRepository) SaveBinding(ctx context.Context, binding resolver.CapabilityBinding, expectedVersion int64) error {
	if r == nil || r.pool == nil {
		return errors.New("repository is not initialized")
	}
	result, err := r.pool.Exec(ctx, `UPDATE capability.capability_binding cb SET binding_mode=$2, priority=$3, status=LOWER($4), contract_version=$5, version=version+1, updated_at=now() FROM capability.capability c WHERE cb.id=$1::uuid AND cb.capability_id=c.capability_id AND c.code=$6 AND cb.version=$7`, binding.ID, binding.BindingMode, binding.Priority, binding.Status, binding.ContractVersion, binding.CapabilityKey, expectedVersion)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("binding %s version conflict or not found", binding.EngineInstanceID)
	}
	return nil
}

func (r *PostgresRepository) ListActiveInstances(ctx context.Context, engineID string) ([]resolver.EngineInstance, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("repository is not initialized")
	}
	rows, err := r.pool.Query(ctx, `
		SELECT engine_instance_id, engine_id, region, environment, status
		FROM topology.engine_instance
		WHERE engine_id = $1::uuid AND UPPER(status) = 'ACTIVE'
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
