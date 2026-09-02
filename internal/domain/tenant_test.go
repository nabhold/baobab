package domain

import "testing"

func TestRegisterTenantValidate(t *testing.T) {
	valid := RegisterTenant{LegalEntityID: "zuribeans_za", TenantID: "zuribeans_za", DisplayName: "Zuri Beans", IsolationStrategy: "schema_per_tenant", ResidencyRegion: "af-south-1", RequestedProducts: []string{"baobab_trade"}}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid command rejected: %v", err)
	}
	valid.RequestedProducts = []string{"baobab_trade", "baobab_trade"}
	if err := valid.Validate(); err == nil {
		t.Fatal("duplicate product accepted")
	}
}

func TestResolveContextRequestValidate(t *testing.T) {
	valid := ResolveContextRequest{TenantID: "zuribeans_za", ProductID: "baobab_trade"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid resolve request rejected: %v", err)
	}
	invalid := ResolveContextRequest{TenantID: "bad id"}
	if err := invalid.Validate(); err == nil {
		t.Fatal("invalid tenant id accepted")
	}
}

func TestEntitlementQueryValidate(t *testing.T) {
	valid := EntitlementQuery{TenantID: "zuribeans_za", ProductID: "baobab_trade"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid entitlement query rejected: %v", err)
	}
	invalid := EntitlementQuery{TenantID: "bad id", ProductID: "baobab_trade"}
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
	valid := LifecycleAction{TenantID: "zuribeans_za", Action: "suspend"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid lifecycle action rejected: %v", err)
	}
	if err := (LifecycleAction{TenantID: "bad id", Action: "activate"}).Validate(); err == nil {
		t.Fatal("invalid lifecycle action accepted")
	}
}
