package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/nabhold/baobab-cp/internal/domain"
	basestore "github.com/nabhold/baobab-cp/internal/store"
)

// TestTenantLifecycleEndToEnd is the regression test recorded as backlog item 2
// in docs/reconciliation/shared-control-plane-audit.md §2.1: it runs
// ApplyMigrations against a real PostgreSQL instance and then exercises
// RegisterTenant -> GetTenant -> ResolveContext -> UpdateTenantLifecycle
// end-to-end, so a migration/store split of the kind fixed in that commit
// (tables the store queries never actually being created) fails a build
// instead of only failing in production.
//
// Set TEST_DATABASE_URL to run it; it is skipped otherwise. It requires a
// real PostgreSQL instance, not a mock or SQLite substitute, for the same
// reason postgres_integration_test.go does.
func TestTenantLifecycleEndToEnd(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping PostgreSQL integration test")
	}
	ctx := context.Background()

	store, err := Open(ctx, url)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()
	if err := store.ApplyMigrations(ctx); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	tenantID := domain.NewTenantID()
	command := domain.RegisterTenant{
		LegalEntityID:     "THAMANI-GLOBAL",
		TenantID:          tenantID,
		DisplayName:       "Thamani Global (integration test)",
		IsolationStrategy: "schema_per_tenant",
		ResidencyRegion:   "af-south-1",
		RequestedProducts: []string{"baobab-trade"},
	}
	if err := command.Validate(); err != nil {
		t.Fatalf("command should be valid: %v", err)
	}

	metadata := basestore.RequestMetadata{ActorID: "integration-test", ActorType: "workload", CorrelationID: "9f8b6e2a-0000-4000-8000-000000000001"}
	idempotencyKey := "integration-test-" + tenantID

	operation, err := store.RegisterTenant(ctx, idempotencyKey, metadata, command)
	if err != nil {
		t.Fatalf("register tenant: %v", err)
	}
	if operation.TenantID != tenantID {
		t.Fatalf("expected operation for tenant %q, got %q", tenantID, operation.TenantID)
	}

	// Replaying the same idempotency key with the same request must return the
	// original operation, not create a second tenant or error.
	replay, err := store.RegisterTenant(ctx, idempotencyKey, metadata, command)
	if err != nil {
		t.Fatalf("idempotent replay: %v", err)
	}
	if replay.OperationID != operation.OperationID {
		t.Fatalf("expected idempotent replay to return the original operation %q, got %q", operation.OperationID, replay.OperationID)
	}

	tenant, err := store.GetTenant(ctx, tenantID)
	if err != nil {
		t.Fatalf("get tenant: %v", err)
	}
	if tenant.LegalEntityID != "THAMANI-GLOBAL" || tenant.DesiredState != "active" {
		t.Fatalf("unexpected tenant record: %#v", tenant)
	}

	// A freshly registered tenant has desired_state=active but
	// observed_state=pending (see 000019_tenant_lifecycle.sql), so context
	// resolution must fail closed until a provisioner confirms activation.
	if _, err := store.ResolveContext(ctx, metadata, tenantID, "baobab-trade"); err == nil {
		t.Fatal("expected context resolution to fail closed before the tenant is observed active")
	}

	if _, err := store.pool.Exec(ctx, `UPDATE tenants SET observed_state='active' WHERE tenant_id=$1`, tenantID); err != nil {
		t.Fatalf("mark tenant observed active: %v", err)
	}
	if _, err := store.pool.Exec(ctx, `UPDATE product_subscriptions SET status='active' WHERE tenant_id=$1 AND product_id='baobab-trade'`, tenantID); err != nil {
		t.Fatalf("mark product subscription active: %v", err)
	}

	resolved, err := store.ResolveContext(ctx, metadata, tenantID, "baobab-trade")
	if err != nil {
		t.Fatalf("resolve context: %v", err)
	}
	if !resolved.Entitled || resolved.TenantID != tenantID || resolved.EntityID != "THAMANI-GLOBAL" {
		t.Fatalf("unexpected resolved context: %#v", resolved)
	}

	if err := store.UpdateTenantLifecycle(ctx, tenantID, domain.LifecycleSuspended); err != nil {
		t.Fatalf("suspend tenant: %v", err)
	}
	tenant, err = store.GetTenant(ctx, tenantID)
	if err != nil {
		t.Fatalf("get tenant after suspend: %v", err)
	}
	if tenant.DesiredState != string(domain.LifecycleSuspended) {
		t.Fatalf("expected tenant to be suspended, got desired_state=%q", tenant.DesiredState)
	}
}
