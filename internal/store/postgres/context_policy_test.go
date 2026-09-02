package postgres

import "testing"

func TestContextPolicyDecision(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		found              bool
		desired            string
		observed           string
		subscription       string
		wantResult         string
		wantPolicyDecision string
	}{
		{name: "active and entitled", found: true, desired: "active", observed: "active", subscription: "active", wantResult: "allowed", wantPolicyDecision: "context_allowed"},
		{name: "unknown tenant", found: false, wantResult: "denied", wantPolicyDecision: "tenant_unknown"},
		{name: "suspended tenant", found: true, desired: "suspended", observed: "active", subscription: "active", wantResult: "denied", wantPolicyDecision: "tenant_not_active"},
		{name: "unreconciled tenant", found: true, desired: "active", observed: "pending", subscription: "active", wantResult: "denied", wantPolicyDecision: "tenant_not_active"},
		{name: "missing entitlement", found: true, desired: "active", observed: "active", wantResult: "denied", wantPolicyDecision: "product_not_entitled"},
		{name: "requested entitlement", found: true, desired: "active", observed: "active", subscription: "requested", wantResult: "denied", wantPolicyDecision: "product_not_entitled"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, policyDecision := contextPolicyDecision(test.found, test.desired, test.observed, test.subscription)
			if result != test.wantResult || policyDecision != test.wantPolicyDecision {
				t.Fatalf("got (%q,%q), want (%q,%q)", result, policyDecision, test.wantResult, test.wantPolicyDecision)
			}
		})
	}
}
