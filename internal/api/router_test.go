package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nabhold/baobab-cp/internal/auth"
	"github.com/nabhold/baobab-cp/internal/domain"
	"github.com/nabhold/baobab-cp/internal/store"
)

type fakeStore struct {
	calls           int
	contextCalls    int
	metadata        store.RequestMetadata
	resolvedTenant  string
	resolvedProduct string
	contextResult   domain.ResolvedContext
	contextErr      error
}

func (f *fakeStore) Ping(context.Context) error { return nil }
func (f *fakeStore) RegisterTenant(_ context.Context, _ string, metadata store.RequestMetadata, command domain.RegisterTenant) (domain.Operation, error) {
	f.calls++
	f.metadata = metadata
	return domain.Operation{OperationID: "7c8f131b-d8ba-4d89-b60b-a187d3944074", TenantID: command.TenantID, State: "accepted", Revision: 1}, nil
}
func (f *fakeStore) ResolveContext(_ context.Context, metadata store.RequestMetadata, tenantID, productID string) (domain.ResolvedContext, error) {
	f.contextCalls++
	f.metadata = metadata
	f.resolvedTenant = tenantID
	f.resolvedProduct = productID
	return f.contextResult, f.contextErr
}

type fakeVerifier struct {
	principal auth.Principal
	err       error
}

func (f fakeVerifier) Verify(context.Context, string) (auth.Principal, error) {
	return f.principal, f.err
}

func adminPrincipal() auth.Principal {
	return auth.Principal{Subject: "admin-123", ActorType: "human", TokenID: "token-123", Scopes: map[string]struct{}{"tenant:write": {}}}
}
func workloadPrincipal() auth.Principal {
	return auth.Principal{Subject: "baobab-trade", ActorType: "workload", TenantID: "tn_01k4m7x9q2v6c8r3d5f1h0j4", ClientID: "baobab-trade", TokenID: "workload-token-123", Scopes: map[string]struct{}{"context:resolve": {}}}
}
func request(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/v1/tenants", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer signed-token")
	req.Header.Set("Idempotency-Key", "tenant-register-0001")
	req.Header.Set("X-Correlation-ID", "7c8f131b-d8ba-4d89-b60b-a187d3944074")
	return req
}

func TestRegisterTenant(t *testing.T) {
	database := &fakeStore{}
	handler := New(Dependencies{Store: database, AdminVerifier: fakeVerifier{principal: adminPrincipal()}})
	body := `{"legal_entity_id":"zuribeans_za","tenant_id":"zuribeans_za","display_name":"Zuri Beans","isolation_strategy":"schema_per_tenant","residency_region":"af-south-1"}`
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request(body))
	if response.Code != http.StatusAccepted {
		t.Fatalf("got status %d: %s", response.Code, response.Body.String())
	}
	if database.calls != 1 {
		t.Fatalf("expected one persistence call, got %d", database.calls)
	}
	if database.metadata.ActorID != "admin-123" || database.metadata.CorrelationID != "7c8f131b-d8ba-4d89-b60b-a187d3944074" {
		t.Fatalf("audit metadata not propagated: %#v", database.metadata)
	}
}

func TestRegisterTenantRejectsInvalidToken(t *testing.T) {
	handler := New(Dependencies{Store: &fakeStore{}, AdminVerifier: fakeVerifier{err: errors.New("invalid")}})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request(`{}`))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d", response.Code)
	}
}

func TestRegisterTenantRequiresHumanWriteScope(t *testing.T) {
	principal := adminPrincipal()
	principal.ActorType = "workload"
	delete(principal.Scopes, "tenant:write")
	handler := New(Dependencies{Store: &fakeStore{}, AdminVerifier: fakeVerifier{principal: principal}})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request(`{}`))
	if response.Code != http.StatusForbidden {
		t.Fatalf("got status %d", response.Code)
	}
}

func TestInvalidCorrelationIDFailsBeforePersistence(t *testing.T) {
	database := &fakeStore{}
	handler := New(Dependencies{Store: database, AdminVerifier: fakeVerifier{principal: adminPrincipal()}})
	req := request(`{}`)
	req.Header.Set("X-Correlation-ID", "not-a-uuid")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("got status %d", response.Code)
	}
	if database.calls != 0 {
		t.Fatal("invalid request reached persistence")
	}
}

func contextRequest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/v1/context/resolve", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer signed-workload-token")
	req.Header.Set("X-Correlation-ID", "7c8f131b-d8ba-4d89-b60b-a187d3944074")
	return req
}

func TestResolveContextReturnsCanonicalEntitledContext(t *testing.T) {
	resolvedAt := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	database := &fakeStore{contextResult: domain.ResolvedContext{
		TenantID: "tn_01k4m7x9q2v6c8r3d5f1h0j4", EntityID: "ZURIBEANS", LifecycleStatus: "active",
		ProductID: "baobab-trade", Entitled: true, CacheTTLSeconds: 15, ResolvedAt: resolvedAt,
		CorrelationID: "7c8f131b-d8ba-4d89-b60b-a187d3944074",
	}}
	handler := New(Dependencies{Store: database, WorkloadVerifier: fakeVerifier{principal: workloadPrincipal()}})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, contextRequest(`{"product_id":"baobab-trade"}`))
	if response.Code != http.StatusOK {
		t.Fatalf("got status %d: %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Cache-Control") != "private, max-age=15" {
		t.Fatalf("unexpected cache policy %q", response.Header().Get("Cache-Control"))
	}
	if database.contextCalls != 1 || database.resolvedTenant != workloadPrincipal().TenantID || database.resolvedProduct != "baobab-trade" {
		t.Fatalf("resolution did not use authenticated tenant context: %#v", database)
	}
	if database.metadata.ClientID != "baobab-trade" || database.metadata.TokenID != "workload-token-123" || database.metadata.CorrelationID == "" {
		t.Fatalf("audit metadata not propagated: %#v", database.metadata)
	}
	var got domain.ResolvedContext
	if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.EntityID != "ZURIBEANS" || !got.Entitled || got.ProductID != "baobab-trade" {
		t.Fatalf("unexpected context response: %#v", got)
	}
}

func TestResolveContextFailsClosed(t *testing.T) {
	tests := []struct {
		name      string
		principal auth.Principal
		storeErr  error
		want      int
	}{
		{name: "human actor", principal: adminPrincipal(), want: http.StatusForbidden},
		{name: "missing scope", principal: func() auth.Principal { p := workloadPrincipal(); p.Scopes = map[string]struct{}{}; return p }(), want: http.StatusForbidden},
		{name: "missing tenant", principal: func() auth.Principal { p := workloadPrincipal(); p.TenantID = ""; return p }(), want: http.StatusForbidden},
		{name: "missing service identity", principal: func() auth.Principal { p := workloadPrincipal(); p.ClientID = ""; return p }(), want: http.StatusForbidden},
		{name: "not entitled", principal: workloadPrincipal(), storeErr: store.ErrContextDenied, want: http.StatusForbidden},
		{name: "store unavailable", principal: workloadPrincipal(), storeErr: errors.New("database unavailable"), want: http.StatusServiceUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			database := &fakeStore{contextErr: test.storeErr}
			handler := New(Dependencies{Store: database, WorkloadVerifier: fakeVerifier{principal: test.principal}})
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, contextRequest(`{"product_id":"baobab-trade"}`))
			if response.Code != test.want {
				t.Fatalf("got status %d: %s", response.Code, response.Body.String())
			}
			if test.want == http.StatusForbidden && strings.Contains(response.Body.String(), "database") {
				t.Fatal("denial leaked internal state")
			}
		})
	}
}

func TestResolveContextRejectsInvalidContractJSON(t *testing.T) {
	for _, body := range []string{`{}`, `{"product_id":"ERP"}`, `{"product_id":"baobab-trade","tenant_id":"attacker"}`, `{"product_id":"baobab-trade"}{}`} {
		database := &fakeStore{}
		handler := New(Dependencies{Store: database, WorkloadVerifier: fakeVerifier{principal: workloadPrincipal()}})
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, contextRequest(body))
		if response.Code != http.StatusBadRequest {
			t.Fatalf("body %q got status %d: %s", body, response.Code, response.Body.String())
		}
		if database.contextCalls != 0 {
			t.Fatalf("body %q reached persistence", body)
		}
	}
}

func TestProblemResponseUsesCanonicalContract(t *testing.T) {
	handler := New(Dependencies{Store: &fakeStore{}, WorkloadVerifier: fakeVerifier{err: errors.New("invalid")}})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, contextRequest(`{"product_id":"baobab-trade"}`))
	if response.Code != http.StatusUnauthorized || response.Header().Get("Content-Type") != "application/problem+json" {
		t.Fatalf("unexpected problem response: %d %q", response.Code, response.Header().Get("Content-Type"))
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"type", "title", "status", "detail", "code", "correlation_id", "retryable"} {
		if _, ok := body[field]; !ok {
			t.Fatalf("canonical problem field %q is missing: %#v", field, body)
		}
	}
}
