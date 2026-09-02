package resolver

import (
	"context"
	"testing"
)

func TestCapabilityResolverResolveUsesHighestPriorityBinding(t *testing.T) {
	resolver := CapabilityResolverImpl{}
	query := CapabilityResolutionQuery{
		CapabilityKey: "baobab_trade",
		Context: Context{
			TenantID:      "tenant-123",
			LegalEntityID: "legal-456",
			MarketID:      "market-789",
			CountryCode:   "ZA",
			CurrencyCode:  "ZAR",
			Locale:        "en-ZA",
		},
		Bindings: []CapabilityBinding{
			{CapabilityKey: "baobab_trade", EngineID: "engine-1", EngineInstanceID: "instance-1", BindingMode: "SECONDARY", Status: "ACTIVE", Priority: 10, ContractVersion: "v1"},
			{CapabilityKey: "baobab_trade", EngineID: "engine-2", EngineInstanceID: "instance-2", BindingMode: "PRIMARY", Status: "ACTIVE", Priority: 100, ContractVersion: "v1"},
		},
	}

	resolved, err := resolver.Resolve(context.Background(), query)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if resolved.BindingMode != "PRIMARY" {
		t.Fatalf("expected PRIMARY binding, got %q", resolved.BindingMode)
	}
	if resolved.EngineID != "engine-2" {
		t.Fatalf("expected engine-2, got %q", resolved.EngineID)
	}
}

func TestCapabilityResolverRejectsUnknownCapability(t *testing.T) {
	resolver := CapabilityResolverImpl{}
	_, err := resolver.Resolve(context.Background(), CapabilityResolutionQuery{
		CapabilityKey: "missing_capability",
		Context: Context{TenantID: "tenant-123"},
		Bindings:     nil,
	})
	if err == nil {
		t.Fatal("expected missing capability error")
	}
}
