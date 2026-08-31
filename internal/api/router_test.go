package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nabhold/baobab-cp/internal/auth"
	"github.com/nabhold/baobab-cp/internal/domain"
	"github.com/nabhold/baobab-cp/internal/store"
)

type fakeStore struct {
	calls    int
	metadata store.RequestMetadata
}

func (f *fakeStore) Ping(context.Context) error { return nil }
func (f *fakeStore) RegisterTenant(_ context.Context, _ string, metadata store.RequestMetadata, command domain.RegisterTenant) (domain.Operation, error) {
	f.calls++
	f.metadata = metadata
	return domain.Operation{OperationID: "7c8f131b-d8ba-4d89-b60b-a187d3944074", TenantID: command.TenantID, State: "accepted", Revision: 1}, nil
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
