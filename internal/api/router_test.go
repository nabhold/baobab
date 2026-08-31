package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nabhold/baobab-cp/internal/domain"
)

type fakeStore struct{ calls int }

func (f *fakeStore) Ping(context.Context) error { return nil }
func (f *fakeStore) RegisterTenant(_ context.Context, _ string, command domain.RegisterTenant) (domain.Operation, error) {
	f.calls++
	return domain.Operation{OperationID: "7c8f131b-d8ba-4d89-b60b-a187d3944074", TenantID: command.TenantID, State: "accepted", Revision: 1}, nil
}

func TestRegisterTenant(t *testing.T) {
	store := &fakeStore{}
	handler := New(store, strings.Repeat("a", 32))
	body := `{"legal_entity_id":"zuribeans_za","tenant_id":"zuribeans_za","display_name":"Zuri Beans","isolation_strategy":"schema_per_tenant","residency_region":"af-south-1"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/tenants", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+strings.Repeat("a", 32))
	req.Header.Set("Idempotency-Key", "tenant-register-0001")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusAccepted {
		t.Fatalf("got status %d: %s", response.Code, response.Body.String())
	}
	if store.calls != 1 {
		t.Fatalf("expected one persistence call, got %d", store.calls)
	}
}

func TestRegisterTenantRequiresAuthorization(t *testing.T) {
	handler := New(&fakeStore{}, strings.Repeat("a", 32))
	req := httptest.NewRequest(http.MethodPost, "/v1/tenants", strings.NewReader(`{}`))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d", response.Code)
	}
}
