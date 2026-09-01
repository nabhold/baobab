package domain

import "testing"

func TestResolveContextValidate(t *testing.T) {
	t.Parallel()
	for _, product := range []string{"baobab-trade", "baobab_erp", "erp"} {
		if err := (ResolveContext{ProductID: product}).Validate(); err != nil {
			t.Fatalf("valid product %q rejected: %v", product, err)
		}
	}
	for _, product := range []string{"", "ERP", "baobab--erp", "baobab erp"} {
		if err := (ResolveContext{ProductID: product}).Validate(); err == nil {
			t.Fatalf("invalid product %q accepted", product)
		}
	}
}
