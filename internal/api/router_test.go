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
	resolveErr      error
	resolved        domain.ResolvedContext
	resolvedCalls   int
	tenant          domain.Tenant
	tenantErr       error
	entitlement     domain.Entitlement
	entitlementErr  error
	entitlementCall int
}

func (f *fakeStore) Ping(context.Context) error { return nil }
func (f *fakeStore) RegisterTenant(_ context.Context, _ string, metadata store.RequestMetadata, command domain.RegisterTenant) (domain.Operation, error) {
	f.calls++
	f.metadata = metadata
	return domain.Operation{OperationID: "7c8f131b-d8ba-4d89-b60b-a187d3944074", TenantID: command.TenantID, State: "accepted", Revision: 1}, nil
}
func (f *fakeStore) ResolveContext(_ context.Context, req domain.ResolveContextRequest) (domain.ResolvedContext, error) {
	f.resolvedCalls++
	if f.resolveErr != nil {
		return domain.ResolvedContext{}, f.resolveErr
	}
	return domain.ResolvedContext{TenantID: req.TenantID, EntityID: req.TenantID, LifecycleStatus: "active", ProductID: req.ProductID, ProductEntitlement: "active"}, nil
}
func (f *fakeStore) GetTenant(_ context.Context, tenantID string) (domain.Tenant, error) {
	if f.tenantErr != nil {
		return domain.Tenant{}, f.tenantErr
	}
	if f.tenant.TenantID == "" {
		f.tenant = domain.Tenant{TenantID: tenantID, LegalEntityID: tenantID, DisplayName: "Zuri Beans", IsolationStrategy: "schema_per_tenant", ResidencyRegion: "af-south-1", DesiredState: "active", ObservedState: "ready", Revision: 1}
	}
	return f.tenant, nil
}
func (f *fakeStore) GetEntitlement(_ context.Context, tenantID, productID string) (domain.Entitlement, error) {
	f.entitlementCall++
	if f.entitlementErr != nil {
		return domain.Entitlement{}, f.entitlementErr
	}
	if f.entitlement.TenantID == "" {
		f.entitlement = domain.Entitlement{TenantID: tenantID, ProductID: productID, Status: "active", Tier: "standard"}
	}
	return f.entitlement, nil
}
func (f *fakeStore) UpdateTenantLifecycle(_ context.Context, tenantID string, next domain.LifecycleStatus) error {
	if !next.Valid() {
		return domain.NotFoundError("invalid tenant lifecycle status")
	}
	return nil
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

func TestResolveContext(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, strings.Repeat("a", 32))
	body := `{"tenant_id":"zuribeans_za","product_id":"baobab_trade"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/context/resolve", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+strings.Repeat("a", 32))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("got status %d: %s", response.Code, response.Body.String())
	}
	if store.resolvedCalls != 1 {
		t.Fatalf("expected one resolve call, got %d", store.resolvedCalls)
	}
}

func TestResolveContextFailsClosed(t *testing.T) {
	store := &fakeStore{resolveErr: storeErrNotFound}
	handler := New(store, strings.Repeat("a", 32))
	req := httptest.NewRequest(http.MethodPost, "/v1/context/resolve", strings.NewReader(`{"tenant_id":"zuribeans_za"}`))
	req.Header.Set("Authorization", "Bearer "+strings.Repeat("a", 32))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusForbidden {
		t.Fatalf("got status %d: %s", response.Code, response.Body.String())
	}
}

func TestGetTenant(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, strings.Repeat("a", 32))
	req := httptest.NewRequest(http.MethodGet, "/v1/tenants/zuribeans_za", nil)
	req.Header.Set("Authorization", "Bearer "+strings.Repeat("a", 32))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("got status %d: %s", response.Code, response.Body.String())
	}
	if response.Body.String() == "" {
		t.Fatal("empty tenant response")
	}
}

func TestGetEntitlement(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, strings.Repeat("a", 32))
	req := httptest.NewRequest(http.MethodGet, "/v1/entitlements?tenantId=zuribeans_za&productId=baobab_trade", nil)
	req.Header.Set("Authorization", "Bearer "+strings.Repeat("a", 32))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("got status %d: %s", response.Code, response.Body.String())
	}
	if store.entitlementCall != 1 {
		t.Fatalf("expected one entitlement lookup, got %d", store.entitlementCall)
	}
}

func TestLifecycleActionEndpoint(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, strings.Repeat("a", 32))
	for _, tc := range []struct {
		path   string
		method string
	}{
		{path: "/v1/tenants/zuribeans_za/suspend", method: http.MethodPost},
		{path: "/v1/tenants/zuribeans_za/activate", method: http.MethodPost},
		{path: "/v1/tenants/zuribeans_za/decommission", method: http.MethodPost},
	} {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		req.Header.Set("Authorization", "Bearer "+strings.Repeat("a", 32))
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		if response.Code != http.StatusOK {
			t.Fatalf("path %s got status %d: %s", tc.path, response.Code, response.Body.String())
		}
	}
}

var storeErrNotFound = domain.NotFoundError("tenant not found")
