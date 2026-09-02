package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nabhold/baobab-cp/internal/auth"
	"github.com/nabhold/baobab-cp/internal/domain"
	"github.com/nabhold/baobab-cp/internal/store"
)

type correlationKey struct{}

type Dependencies struct {
	Store            store.TenantStore
	AdminVerifier    auth.TokenVerifier
	WorkloadVerifier auth.TokenVerifier
}
type API struct {
	store            store.TenantStore
	adminVerifier    auth.TokenVerifier
	workloadVerifier auth.TokenVerifier
}

func New(dependencies Dependencies) http.Handler {
	a := &API{store: dependencies.Store, adminVerifier: dependencies.AdminVerifier, workloadVerifier: dependencies.WorkloadVerifier}
	r := chi.NewRouter()
	r.Use(a.securityHeaders, a.correlation, a.requestLog)
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Get("/readyz", a.ready)
	r.With(a.authorize).Post("/v1/tenants", a.register)
	r.With(a.authorize).Get("/v1/tenants/{tenantID}", a.getTenant)
	r.With(a.authorize).Post("/v1/tenants/{tenantID}/suspend", a.tenantLifecycleAction("suspend"))
	r.With(a.authorize).Post("/v1/tenants/{tenantID}/activate", a.tenantLifecycleAction("activate"))
	r.With(a.authorize).Post("/v1/tenants/{tenantID}/decommission", a.tenantLifecycleAction("decommission"))
	r.With(a.authorize).Get("/v1/entitlements", a.getEntitlement)
	r.With(a.authorize).Post("/v1/context/resolve", a.resolveContext)
	return r
}

func (a *API) authorize(verifier auth.TokenVerifier, actorType, requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				problem(w, r, http.StatusUnauthorized, "AUTH_TOKEN_REQUIRED", "a bearer token is required", false)
				return
			}
			if verifier == nil {
				problem(w, r, http.StatusServiceUnavailable, "AUTH_VERIFIER_UNAVAILABLE", "authentication is temporarily unavailable", true)
				return
			}
			principal, err := verifier.Verify(r.Context(), raw)
			if err != nil {
				slog.WarnContext(r.Context(), "authentication denied", "reason", "invalid_token", "correlation_id", correlationID(r))
				problem(w, r, http.StatusUnauthorized, "AUTH_TOKEN_INVALID", "the bearer token is invalid", false)
				return
			}
			workloadContextMissing := actorType == "workload" && (principal.TenantID == "" || principal.ClientID == "")
			if principal.ActorType != actorType || !principal.HasScope(requiredScope) || workloadContextMissing {
				slog.WarnContext(r.Context(), "authorization denied", "actor_id", principal.Subject, "actor_type", principal.ActorType, "client_id", principal.ClientID, "tenant_id", principal.TenantID, "required_scope", requiredScope, "correlation_id", correlationID(r))
				problem(w, r, http.StatusForbidden, "AUTHORIZATION_DENIED", "the authenticated principal lacks required authority", false)
				return
			}
			*r = *r.WithContext(auth.WithPrincipal(r.Context(), principal))
			next.ServeHTTP(w, r)
		})
	}
}

func bearerToken(header string) (string, bool) {
	scheme, value, ok := strings.Cut(strings.TrimSpace(header), " ")
	return value, ok && strings.EqualFold(scheme, "Bearer") && value != ""
}

func (a *API) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func (a *API) correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Correlation-ID")
		if id != "" && !validUUID(id) {
			id = newUUID()
			w.Header().Set("X-Correlation-ID", id)
			r = r.WithContext(context.WithValue(r.Context(), correlationKey{}, id))
			problem(w, r, http.StatusBadRequest, "INVALID_CORRELATION_ID", "X-Correlation-ID must be a UUID", false)
			return
		}
		if id == "" {
			id = newUUID()
		}
		w.Header().Set("X-Correlation-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), correlationKey{}, id)))
	})
}

func (a *API) requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		response := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(response, r)
		principal, _ := auth.PrincipalFromContext(r.Context())
		slog.InfoContext(r.Context(), "request completed", "method", r.Method, "path", r.URL.Path, "status", response.status, "actor_id", principal.Subject, "actor_type", principal.ActorType, "client_id", principal.ClientID, "tenant_id", principal.TenantID, "correlation_id", correlationID(r), "duration_ms", time.Since(started).Milliseconds())
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (a *API) ready(w http.ResponseWriter, r *http.Request) {
	if err := a.store.Ping(r.Context()); err != nil {
		problem(w, r, http.StatusServiceUnavailable, "SERVICE_NOT_READY", "database unavailable", true)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (a *API) register(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if len(key) < 16 || len(key) > 128 {
		problem(w, r, http.StatusBadRequest, "INVALID_IDEMPOTENCY_KEY", "Idempotency-Key must contain 16 to 128 characters", false)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var command domain.RegisterTenant
	if err := dec.Decode(&command); err != nil {
		problem(w, r, http.StatusBadRequest, "INVALID_REQUEST", "request body is not valid contract JSON", false)
		return
	}
	if err := command.Validate(); err != nil {
		problem(w, r, http.StatusBadRequest, "VALIDATION_FAILED", err.Error(), false)
		return
	}
	principal, _ := auth.PrincipalFromContext(r.Context())
	metadata := requestMetadata(r, principal)
	operation, err := a.store.RegisterTenant(r.Context(), key, metadata, command)
	if errors.Is(err, store.ErrIdempotencyConflict) {
		problem(w, r, http.StatusConflict, "IDEMPOTENCY_KEY_REUSED", "the idempotency key was used for a different request", false)
		return
	}
	if err != nil {
		problem(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "tenant registration could not be persisted", true)
		return
	}
	w.Header().Set("Location", "/v1/operations/"+operation.OperationID)
	writeJSON(w, http.StatusAccepted, operation)
}

func (a *API) resolveContext(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req domain.ResolveContextRequest
	if err := dec.Decode(&req); err != nil {
		problem(w, 400, "invalid_request", "request body is not valid contract JSON")
		return
	}
	if err := req.Validate(); err != nil {
		problem(w, 422, "validation_failed", err.Error())
		return
	}
	ctx, err := a.store.ResolveContext(r.Context(), req)
	if err != nil {
		var notFound domain.NotFoundError
		if errors.As(err, &notFound) {
			problem(w, 403, "tenant_context_unresolved", err.Error())
			return
		}
		problem(w, 500, "internal_error", "tenant context could not be resolved")
		return
	}
	writeJSON(w, 200, ctx)
}

func (a *API) getTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	if !domain.ValidResource(tenantID) {
		problem(w, 400, "invalid_tenant_id", "tenant_id is invalid")
		return
	}
	tenant, err := a.store.GetTenant(r.Context(), tenantID)
	if err != nil {
		var notFound domain.NotFoundError
		if errors.As(err, &notFound) {
			problem(w, 404, "tenant_not_found", err.Error())
			return
		}
		problem(w, 500, "internal_error", "tenant lookup failed")
		return
	}
	writeJSON(w, 200, tenant)
}

func (a *API) getEntitlement(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenantId")
	productID := r.URL.Query().Get("productId")
	q := domain.EntitlementQuery{TenantID: tenantID, ProductID: productID}
	if err := q.Validate(); err != nil {
		problem(w, 422, "validation_failed", err.Error())
		return
	}
	ent, err := a.store.GetEntitlement(r.Context(), tenantID, productID)
	if err != nil {
		var notFound domain.NotFoundError
		if errors.As(err, &notFound) {
			problem(w, 404, "entitlement_not_found", err.Error())
			return
		}
		problem(w, 500, "internal_error", "entitlement lookup failed")
		return
	}
	writeJSON(w, 200, ent)
}

func (a *API) tenantLifecycleAction(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := chi.URLParam(r, "tenantID")
		cmd := domain.LifecycleAction{TenantID: tenantID, Action: action}
		if err := cmd.Validate(); err != nil {
			problem(w, 422, "validation_failed", err.Error())
			return
		}
		var next domain.LifecycleStatus
		switch action {
		case "activate":
			next = domain.LifecycleActive
		case "suspend":
			next = domain.LifecycleSuspended
		case "decommission":
			next = domain.LifecycleDecommissioned
		}
		if err := a.store.UpdateTenantLifecycle(r.Context(), tenantID, next); err != nil {
			var notFound domain.NotFoundError
			if errors.As(err, &notFound) {
				problem(w, 404, "tenant_not_found", err.Error())
				return
			}
			problem(w, 500, "internal_error", "tenant lifecycle update failed")
			return
		}
		writeJSON(w, 200, map[string]string{"tenant_id": tenantID, "status": string(next)})
	}
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func problem(w http.ResponseWriter, r *http.Request, status int, code, detail string, retryable bool) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"type": "https://docs.nabhold.com/problems/" + strings.ToLower(code), "title": http.StatusText(status), "status": status, "detail": detail, "code": code, "correlation_id": correlationID(r), "retryable": retryable})
}

func requestMetadata(r *http.Request, principal auth.Principal) store.RequestMetadata {
	return store.RequestMetadata{ActorID: principal.Subject, ActorType: principal.ActorType, ClientID: principal.ClientID, TokenID: principal.TokenID, CorrelationID: correlationID(r)}
}
