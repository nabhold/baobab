package domain

import "testing"

func TestRegisterTenantValidate(t *testing.T) {
	valid := RegisterTenant{LegalEntityID: "THAMANI-GLOBAL", TenantID: NewTenantID(), DisplayName: "Zuri Beans", IsolationStrategy: "schema_per_tenant", ResidencyRegion: "af-south-1", RequestedProducts: []string{"baobab-trade"}}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid command rejected: %v", err)
	}
	valid.RequestedProducts = []string{"baobab-trade", "baobab-trade"}
	if err := valid.Validate(); err == nil {
		t.Fatal("duplicate product accepted")
	}
}

func TestRegisterTenantValidateAcceptsLegacyLegalEntityAlias(t *testing.T) {
	// ADR-0003 §2.3: Control Plane v1 accepts the former lowercase alias at
	// this input boundary during the documented compatibility window.
	valid := RegisterTenant{LegalEntityID: "zuribeans_za", TenantID: NewTenantID(), DisplayName: "Zuri Beans", IsolationStrategy: "schema_per_tenant", ResidencyRegion: "af-south-1"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("legacy legal_entity_id alias rejected: %v", err)
	}
}

func TestRegisterTenantValidateRejectsMissingTenantID(t *testing.T) {
	// tenant_id is minted by the handler, not the caller; Validate still
	// guards against a programming error that skips minting it.
	invalid := RegisterTenant{LegalEntityID: "THAMANI-GLOBAL", DisplayName: "Zuri Beans", IsolationStrategy: "schema_per_tenant", ResidencyRegion: "af-south-1"}
	if err := invalid.Validate(); err == nil {
		t.Fatal("missing tenant_id accepted")
	}
}

func TestEntitlementQueryValidate(t *testing.T) {
	valid := EntitlementQuery{TenantID: NewTenantID(), ProductID: "baobab-trade"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid entitlement query rejected: %v", err)
	}
	invalid := EntitlementQuery{TenantID: "bad id", ProductID: "baobab-trade"}
	if err := invalid.Validate(); err == nil {
		t.Fatal("invalid tenant id accepted")
	}
}

func TestLifecycleStatusValidate(t *testing.T) {
	for _, status := range []LifecycleStatus{LifecycleProvisioning, LifecycleActive, LifecycleSuspended, LifecycleDecommissioning, LifecycleDecommissioned} {
		if !status.Valid() {
			t.Fatalf("status %q should be valid", status)
		}
	}
	if LifecycleStatus("unknown").Valid() {
		t.Fatal("unknown lifecycle status should be rejected")
	}
}

func TestLifecycleTransition(t *testing.T) {
	if next, ok := TransitionLifecycle(LifecycleProvisioning, LifecycleActive); !ok || next != LifecycleActive {
		t.Fatal("expected provisioning to become active")
	}
	if next, ok := TransitionLifecycle(LifecycleActive, LifecycleSuspended); !ok || next != LifecycleSuspended {
		t.Fatal("expected active to become suspended")
	}
	if _, ok := TransitionLifecycle(LifecycleDecommissioned, LifecycleActive); ok {
		t.Fatal("decommissioned tenant should not be re-activated without explicit reset")
	}
}

func TestLifecycleActionValidate(t *testing.T) {
	valid := LifecycleAction{TenantID: NewTenantID(), Action: "suspend"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid lifecycle action rejected: %v", err)
	}
	if err := (LifecycleAction{TenantID: "bad id", Action: "activate"}).Validate(); err == nil {
		t.Fatal("invalid lifecycle action accepted")
	}
}
