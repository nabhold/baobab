package resolver

import (
	"context"
	"testing"
)

func TestTopologyResolverSelectsActiveInstance(t *testing.T) {
	resolver := TopologyResolverImpl{}
	query := TopologyResolutionQuery{
		Context: Context{TenantID: "tenant-123", MarketID: "market-789", CountryCode: "ZA"},
		EngineInstances: []EngineInstance{
			{ID: "instance-1", EngineID: "engine-1", Region: "af-south-1", Environment: "production", Status: "ACTIVE"},
			{ID: "instance-2", EngineID: "engine-1", Region: "af-south-1", Environment: "staging", Status: "ACTIVE"},
		},
	}

	resolved, err := resolver.Resolve(context.Background(), query)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if resolved.ID != "instance-1" {
		t.Fatalf("expected instance-1, got %q", resolved.ID)
	}
}

func TestTopologyResolverRejectsUnavailableEngine(t *testing.T) {
	resolver := TopologyResolverImpl{}
	_, err := resolver.Resolve(context.Background(), TopologyResolutionQuery{
		Context: Context{TenantID: "tenant-123"},
		EngineInstances: []EngineInstance{{
			ID:          "instance-1",
			EngineID:    "engine-1",
			Region:      "af-south-1",
			Environment: "production",
			Status:      "MAINTENANCE",
		}},
	})
	if err == nil {
		t.Fatal("expected unavailable engine rejection")
	}
}
