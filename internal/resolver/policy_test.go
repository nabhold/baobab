package resolver

import (
	"context"
	"testing"
)

func TestPolicyCheckerAllowsActiveBindingWithContext(t *testing.T) {
	checker := PolicyChecker{}
	decision := checker.Check(context.Background(), Context{TenantID: "tenant-123", MarketID: "market-789", CountryCode: "ZA"}, CapabilityBinding{
		CapabilityKey:    "baobab_trade",
		EngineID:         "engine-1",
		EngineInstanceID: "instance-1",
		BindingMode:      "PRIMARY",
		Status:           "ACTIVE",
	})
	if !decision.Allowed {
		t.Fatalf("expected policy allow, got reason %q", decision.Reason)
	}
}

func TestPolicyCheckerRejectsExpiredOrMissingBinding(t *testing.T) {
	checker := PolicyChecker{}
	decision := checker.Check(context.Background(), Context{TenantID: "tenant-123"}, CapabilityBinding{
		CapabilityKey:    "baobab_trade",
		EngineID:         "engine-1",
		EngineInstanceID: "",
		BindingMode:      "PRIMARY",
		Status:           "ACTIVE",
	})
	if decision.Allowed {
		t.Fatal("expected policy rejection for missing engine instance")
	}
}
