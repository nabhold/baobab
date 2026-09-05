package domain_test

import (
	"testing"
	"time"

	"github.com/nabhold/baobab-cp/internal/contracttest"
	"github.com/nabhold/baobab-cp/internal/domain"
)

// These tests validate payloads this repository's own domain types produce
// against the actual JSON Schemas in a local nabhold/shared checkout, per
// the automation called for in
// docs/reconciliation/shared-control-plane-audit.md §6/§10.5. Set
// SHARED_CONTRACTS_DIR to a checkout of nabhold/shared at (or above) the
// commit pinned in contracts.lock.yaml; they are skipped otherwise.

func TestRegisterTenantMatchesSharedSchema(t *testing.T) {
	dir := contracttest.SharedDir(t)
	schema := contracttest.CompileSchema(t, dir, "control-plane/v1/tenant-registration.schema.json")

	command := domain.RegisterTenant{
		LegalEntityID:     "THAMANI-GLOBAL",
		TenantID:          domain.NewTenantID(), // json:"-": must not appear in the encoded payload
		DisplayName:       "Thamani Global",
		IsolationStrategy: "schema_per_tenant",
		ResidencyRegion:   "af-south-1",
		RequestedProducts: []string{"baobab-trade"},
	}
	if err := command.Validate(); err != nil {
		t.Fatalf("command should be valid: %v", err)
	}
	contracttest.ValidateJSON(t, schema, command)
}

func TestResolvedContextMatchesSharedSchema(t *testing.T) {
	dir := contracttest.SharedDir(t)
	schema := contracttest.CompileSchema(t, dir, "control-plane/v1/context-resolution.schema.json#/$defs/response")

	tier := "standard"
	resolved := domain.ResolvedContext{
		TenantID:        "tn_01k4example",
		EntityID:        "THAMANI-GLOBAL",
		LifecycleStatus: "active",
		ProductID:       "baobab-trade",
		Entitled:        true,
		EntitlementTier: &tier,
		CacheTTLSeconds: 15,
		ResolvedAt:      time.Now().UTC(),
		CorrelationID:   "9f8b6e2a-0000-4000-8000-000000000001",
	}
	contracttest.ValidateJSON(t, schema, resolved)
}
