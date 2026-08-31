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
