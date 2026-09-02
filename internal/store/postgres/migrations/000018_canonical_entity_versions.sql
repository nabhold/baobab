ALTER TABLE registry.canonical_entity
    ADD COLUMN IF NOT EXISTS version bigint NOT NULL DEFAULT 1;

CREATE INDEX IF NOT EXISTS canonical_entity_version_idx
    ON registry.canonical_entity(version);
