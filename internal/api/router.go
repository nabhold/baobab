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
	r.With(a.authorize(a.adminVerifier, "human", "tenant:write")).Post("/v1/tenants", a.register)
	r.With(a.authorize(a.workloadVerifier, "workload", "context:resolve")).Post("/v1/context/resolve", a.resolveContext)
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

func correlationID(r *http.Request) string {
	id, _ := r.Context().Value(correlationKey{}).(string)
	return id
}
func validUUID(value string) bool {
	if len(value) != 36 {
		return false
	}
	for i, c := range value {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
			continue
		}
		if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
			return false
		}
	}
	return true
}
func newUUID() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		panic("cryptographic random source unavailable")
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	encoded := make([]byte, 36)
	hex.Encode(encoded[0:8], value[0:4])
	encoded[8] = '-'
	hex.Encode(encoded[9:13], value[4:6])
	encoded[13] = '-'
	hex.Encode(encoded[14:18], value[6:8])
	encoded[18] = '-'
	hex.Encode(encoded[19:23], value[8:10])
	encoded[23] = '-'
	hex.Encode(encoded[24:36], value[10:16])
	return string(encoded)
}
func writeJSON(w http.ResponseWriter, status int, value any) {
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
