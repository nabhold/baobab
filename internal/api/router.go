package api

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/nabhold/baobab-cp/internal/domain"
	"github.com/nabhold/baobab-cp/internal/store"
	"net/http"
	"strings"
)

type API struct {
	store store.TenantStore
	token string
}

func New(s store.TenantStore, token string) http.Handler {
	a := &API{s, token}
	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, map[string]string{"status": "ok"}) })
	r.Get("/readyz", a.ready)
	r.With(a.authorize).Post("/v1/tenants", a.register)
	return r
}
func (a *API) authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if len(got) != len(a.token) || subtle.ConstantTimeCompare([]byte(got), []byte(a.token)) != 1 {
			problem(w, 401, "unauthorized", "a valid administrative bearer token is required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (a *API) ready(w http.ResponseWriter, r *http.Request) {
	if err := a.store.Ping(r.Context()); err != nil {
		problem(w, 503, "not_ready", "database unavailable")
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
func (a *API) register(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if len(key) < 16 || len(key) > 128 {
		problem(w, 400, "invalid_idempotency_key", "Idempotency-Key must contain 16 to 128 characters")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var c domain.RegisterTenant
	if err := dec.Decode(&c); err != nil {
		problem(w, 400, "invalid_request", "request body is not valid contract JSON")
		return
	}
	if err := c.Validate(); err != nil {
		problem(w, 422, "validation_failed", err.Error())
		return
	}
	op, err := a.store.RegisterTenant(r.Context(), key, c)
	if errors.Is(err, store.ErrIdempotencyConflict) {
		problem(w, 409, "idempotency_conflict", err.Error())
		return
	}
	if err != nil {
		problem(w, 500, "internal_error", "tenant registration could not be persisted")
		return
	}
	w.Header().Set("Location", "/v1/operations/"+op.OperationID)
	writeJSON(w, 202, op)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func problem(w http.ResponseWriter, status int, code, detail string) {
	writeJSON(w, status, map[string]any{"type": "https://docs.nabhold.com/problems/" + code, "title": http.StatusText(status), "status": status, "detail": detail, "code": code})
}
