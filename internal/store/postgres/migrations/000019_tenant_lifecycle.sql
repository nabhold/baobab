-- Renumbered from the previously unregistered 000001_control_plane.up.sql; see
-- docs/reconciliation/shared-control-plane-audit.md §2.1.
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE legal_entities (legal_entity_id varchar(63) PRIMARY KEY,created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE tenants (tenant_id varchar(63) PRIMARY KEY,legal_entity_id varchar(63) NOT NULL REFERENCES legal_entities,display_name varchar(255) NOT NULL,isolation_strategy varchar(32) NOT NULL CHECK(isolation_strategy IN('schema_per_tenant','row_level_security')),residency_region varchar(64) NOT NULL,metadata jsonb NOT NULL DEFAULT '{}',desired_state varchar(32) NOT NULL DEFAULT 'active',observed_state varchar(32) NOT NULL DEFAULT 'pending',revision bigint NOT NULL DEFAULT 1,created_at timestamptz NOT NULL DEFAULT now(),updated_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE product_subscriptions(tenant_id varchar(63) NOT NULL REFERENCES tenants ON DELETE CASCADE,product_id varchar(63) NOT NULL,status varchar(32) NOT NULL DEFAULT 'requested',PRIMARY KEY(tenant_id,product_id));
CREATE TABLE provisioning_operations(operation_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),tenant_id varchar(63) NOT NULL REFERENCES tenants,idempotency_key varchar(128) NOT NULL UNIQUE,request_hash char(64) NOT NULL,state varchar(32) NOT NULL DEFAULT 'accepted',revision integer NOT NULL DEFAULT 1,created_at timestamptz NOT NULL DEFAULT now(),updated_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE outbox_events(event_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),aggregate_id varchar(63) NOT NULL,event_type varchar(128) NOT NULL,payload jsonb NOT NULL,occurred_at timestamptz NOT NULL DEFAULT now(),published_at timestamptz);
CREATE INDEX outbox_unpublished_idx ON outbox_events(occurred_at) WHERE published_at IS NULL;
CREATE TABLE audit_events(audit_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),tenant_id varchar(63),action varchar(128) NOT NULL,payload jsonb NOT NULL,occurred_at timestamptz NOT NULL DEFAULT now());
REVOKE UPDATE,DELETE ON audit_events FROM PUBLIC;
