-- NOTE (implementation-conformance review, see docs/adr/ADR-0005-bcp-db-001-conformance-gap.md):
--
-- BCP-DB-001 (docs/adr/ADR-BCP-001-...md, section 22) specifies this migration
-- against `registry.canonical_entity.owner_tenant_id`, a `registry.canonical_relationship`
-- with `source_entity_id`/`target_entity_id`/`status`/`valid_period`, a `mapping.mapping`
-- table carrying `external_reference_id`/`target_canonical_entity_id`/`scope_id`/
-- `direction`/`cardinality`/`authority`/`confidence`/`valid_period`, and a
-- `mapping.mapping_scope` carrying `digital_estate_id`/`engine_id`/`engine_instance_id`/
-- `operating_region_id`/`country_code`/`currency_code`/`locale`.
--
-- The migrations actually committed ahead of this one (000002, 000003, 000009,
-- 000011) implement a materially smaller schema: `registry.canonical_entity.tenant_id`,
-- a `registry.canonical_relationship` with only `parent_entity_id`/`child_entity_id`
-- (no status or temporal validity), a `mapping.canonical_mapping` table with only
-- `source_entity_id`/`target_entity_id`/`mapping_type`/`status` (no scope, no
-- external-reference linkage, no temporal exclusion), and a `mapping.mapping_scope`
-- with only `tenant_id`/`market_id`/`entity_type`/`scope_type`.
--
-- As originally written, this file referenced columns and a `mapping.mapping`
-- table that do not exist, and `cmd/migrate` failed outright on a fresh database
-- (verified: PostgreSQL 16/17, `column "owner_tenant_id" does not exist`,
-- `relation "mapping.mapping" does not exist`, and ten further errors of the
-- same kind). That is a deployment-blocking defect independent of which schema
-- is ultimately correct, so it is fixed here to match the schema that actually
-- exists. The indexes BCP-DB-001 specifies for the richer, temporally-exclusive
-- `mapping.mapping` / `capability.capability_binding`-style model are deferred,
-- not silently dropped: see the ADR for the two remediation options (elevate
-- 000002-000011 to BCP-DB-001 fidelity, or formally supersede BCP-DB-001's
-- mapping/relationship model with the simplified one below).

CREATE INDEX IF NOT EXISTS canonical_entity_type_status_idx
    ON registry.canonical_entity(entity_type, status);

CREATE INDEX IF NOT EXISTS canonical_entity_tenant_type_idx
    ON registry.canonical_entity(tenant_id, entity_type)
    WHERE tenant_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS canonical_relationship_parent_type_idx
    ON registry.canonical_relationship(parent_entity_id, relationship_type);

CREATE INDEX IF NOT EXISTS canonical_relationship_child_type_idx
    ON registry.canonical_relationship(child_entity_id, relationship_type);

CREATE INDEX IF NOT EXISTS mapping_scope_tenant_market_idx
    ON mapping.mapping_scope(tenant_id, market_id);

CREATE INDEX IF NOT EXISTS mapping_scope_entity_type_idx
    ON mapping.mapping_scope(entity_type, scope_type);

CREATE INDEX IF NOT EXISTS canonical_mapping_source_lookup_idx
    ON mapping.canonical_mapping(source_entity_id, mapping_type, status);

CREATE INDEX IF NOT EXISTS canonical_mapping_target_lookup_idx
    ON mapping.canonical_mapping(target_entity_id, mapping_type, status);

CREATE INDEX IF NOT EXISTS canonical_mapping_status_created_idx
    ON mapping.canonical_mapping(status, created_at DESC);

CREATE INDEX IF NOT EXISTS capability_binding_lookup_idx
    ON capability.capability_binding(capability_id, status, binding_mode);

CREATE INDEX IF NOT EXISTS capability_binding_scope_idx
    ON capability.capability_binding(scope_id);

CREATE INDEX IF NOT EXISTS capability_binding_period_gist
    ON capability.capability_binding USING gist(valid_period);

CREATE INDEX IF NOT EXISTS outbox_pending_idx
    ON messaging.outbox(coalesce(next_attempt_at, occurred_at), occurred_at)
    WHERE published_at IS NULL;

CREATE INDEX IF NOT EXISTS inbox_processing_idx
    ON messaging.inbox(processing_status, received_at);

CREATE INDEX IF NOT EXISTS audit_resource_idx
    ON audit.audit_record(resource_type, resource_id, occurred_at DESC);

CREATE INDEX IF NOT EXISTS audit_tenant_time_idx
    ON audit.audit_record(tenant_id, occurred_at DESC)
    WHERE tenant_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS audit_correlation_idx
    ON audit.audit_record(correlation_id)
    WHERE correlation_id IS NOT NULL;
