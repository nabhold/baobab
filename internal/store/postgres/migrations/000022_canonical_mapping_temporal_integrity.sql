-- Adds the temporal-validity columns and non-overlap invariant the Canonical
-- Mapping Model requires (§25: "No two mutually exclusive authoritative
-- mappings SHOULD overlap for the same canonical subject") and that
-- contracts/control-plane/v1/canonical-mapping.schema.json's mapping
-- definition already requires (effective_from is a required field). Until
-- this migration, mapping.canonical_mapping had no temporal columns at all:
-- domain.Mapping.Validate() required and parsed effective_from/effective_to,
-- but internal/repository/postgres.go silently discarded them before
-- persisting and fabricated EffectiveFrom from created_at when reading rows
-- back. See docs/reconciliation/shared-control-plane-audit.md §2.2.
ALTER TABLE mapping.canonical_mapping
    ADD COLUMN effective_from timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN effective_to timestamptz;

ALTER TABLE mapping.canonical_mapping
    ADD CONSTRAINT canonical_mapping_period_ck
    CHECK (effective_to IS NULL OR effective_to > effective_from);

ALTER TABLE mapping.canonical_mapping
    ADD COLUMN valid_period tstzrange
    GENERATED ALWAYS AS (tstzrange(effective_from, effective_to, '[)')) STORED;

-- A source entity must not have two simultaneously active mappings of the
-- same type with overlapping validity windows - the "ambiguous authoritative
-- mapping" outcome the Mapping Model requires resolution to fail on rather
-- than silently pick a winner. This is a narrower version of BCP-DB-001's
-- mapping_single_authoritative_excl pattern, scoped to the columns this
-- table actually has (no scope_id/resolution_mode/confidence exist here yet
-- - see the audit's §10.4 for that larger, separately-scoped rework).
ALTER TABLE mapping.canonical_mapping
    ADD CONSTRAINT canonical_mapping_source_type_active_excl
    EXCLUDE USING gist (
        source_entity_id WITH =,
        mapping_type WITH =,
        valid_period WITH &&
    )
    WHERE (status = 'active');
