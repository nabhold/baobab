-- messaging.outbox (migration 000015) typed aggregate_id and tenant_id as
-- uuid, but the actual identifier grammar for a tenant-provisioning
-- aggregate is the opaque, Control Plane-minted tn_-prefixed string from
-- contracts/control-plane/v1/domain.schema.json's tenantId definition (see
-- internal/domain/ids.go), not a uuid. A uuid column cannot store a "tn_..."
-- value at all, so no tenant-provisioning event could ever be written to
-- this table as committed. Aggregate identity is polymorphic across bounded
-- contexts (a uuid for a canonical-entity/mapping aggregate, a tn_-prefixed
-- string for a tenant aggregate), so both columns are widened to text
-- rather than narrowed to either shape.
--
-- See docs/reconciliation/shared-control-plane-audit.md §12 (P1 backlog
-- item 6) for the finding this corrects.
ALTER TABLE messaging.outbox
    ALTER COLUMN aggregate_id TYPE text USING aggregate_id::text,
    ALTER COLUMN tenant_id TYPE text USING tenant_id::text;
