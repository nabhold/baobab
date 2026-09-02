CREATE TABLE IF NOT EXISTS capability.capability_binding (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    capability_id uuid NOT NULL REFERENCES capability.capability(capability_id) ON DELETE CASCADE,
    engine_instance_id uuid NOT NULL REFERENCES topology.engine_instance(engine_instance_id) ON DELETE CASCADE,
    scope_id uuid NOT NULL REFERENCES mapping.mapping_scope(mapping_scope_id) ON DELETE CASCADE,
    binding_mode text NOT NULL,
    priority integer NOT NULL DEFAULT 100,
    status text NOT NULL,
    contract_version text NOT NULL,
    effective_from timestamptz NOT NULL,
    effective_to timestamptz,
    valid_period tstzrange GENERATED ALWAYS AS (tstzrange(effective_from, effective_to, '[)')) STORED,
    fallback_binding_id uuid REFERENCES capability.capability_binding(id) ON DELETE SET NULL,
    configuration jsonb NOT NULL DEFAULT '{}'::jsonb,
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid,
    updated_at timestamptz NOT NULL DEFAULT now(),
    updated_by uuid,
    approved_at timestamptz,
    approved_by uuid,
    CHECK (effective_to IS NULL OR effective_to > effective_from),
    CHECK (binding_mode IN ('PRIMARY', 'SECONDARY', 'FALLBACK', 'READ_ONLY', 'MIGRATION_SOURCE', 'MIGRATION_TARGET', 'SHADOW'))
);

ALTER TABLE capability.capability_binding
ADD CONSTRAINT capability_binding_primary_excl
EXCLUDE USING gist (
    capability_id WITH =,
    scope_id WITH =,
    valid_period WITH &&
)
WHERE (
    status = 'ACTIVE'
    AND binding_mode = 'PRIMARY'
);
