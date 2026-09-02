CREATE INDEX IF NOT EXISTS canonical_entity_type_status_idx
    ON registry.canonical_entity(entity_type, status);

CREATE INDEX IF NOT EXISTS canonical_entity_tenant_type_idx
    ON registry.canonical_entity(owner_tenant_id, entity_type)
    WHERE owner_tenant_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS canonical_relationship_source_idx
    ON registry.canonical_relationship(source_entity_id, relationship_type, status);

CREATE INDEX IF NOT EXISTS canonical_relationship_target_idx
    ON registry.canonical_relationship(target_entity_id, relationship_type, status);

CREATE INDEX IF NOT EXISTS canonical_relationship_period_gist
    ON registry.canonical_relationship USING gist(valid_period);

CREATE INDEX IF NOT EXISTS mapping_scope_tenant_market_idx
    ON mapping.mapping_scope(tenant_id, market_id);

CREATE INDEX IF NOT EXISTS mapping_scope_estate_idx
    ON mapping.mapping_scope(digital_estate_id);

CREATE INDEX IF NOT EXISTS mapping_scope_engine_idx
    ON mapping.mapping_scope(engine_id, engine_instance_id);

CREATE INDEX IF NOT EXISTS mapping_scope_geography_idx
    ON mapping.mapping_scope(operating_region_id, country_code);

CREATE INDEX IF NOT EXISTS mapping_scope_currency_locale_idx
    ON mapping.mapping_scope(currency_code, locale);

CREATE INDEX IF NOT EXISTS mapping_canonical_lookup_idx
    ON mapping.mapping(canonical_entity_id, mapping_type, status);

CREATE INDEX IF NOT EXISTS mapping_external_lookup_idx
    ON mapping.mapping(external_reference_id, status)
    WHERE external_reference_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS mapping_target_canonical_idx
    ON mapping.mapping(target_canonical_entity_id, status)
    WHERE target_canonical_entity_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS mapping_scope_lookup_idx
    ON mapping.mapping(scope_id, mapping_type, status);

CREATE INDEX IF NOT EXISTS mapping_valid_period_gist
    ON mapping.mapping USING gist(valid_period);

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
