package api

import (
	"encoding/json"
	"net/http"

	"github.com/nabhold/baobab-cp/internal/domain"
	"github.com/nabhold/baobab-cp/internal/resolver"
	"github.com/nabhold/baobab-cp/internal/service"
)

// ResolverHandler exposes the composed resolution service over HTTP.
type ResolverHandler struct {
	Service service.ResolutionService
}

type resolverRequest struct {
	TenantID        string                       `json:"tenant_id"`
	Context         resolver.Context             `json:"context"`
	Mappings        []domain.Mapping             `json:"mappings"`
	Bindings        []resolver.CapabilityBinding `json:"bindings"`
	EngineInstances []resolver.EngineInstance    `json:"engine_instances"`
}

func (h ResolverHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		problem(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "POST only", false)
		return
	}

	var req resolverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem(w, r, http.StatusBadRequest, "INVALID_REQUEST", err.Error(), false)
		return
	}
	if req.TenantID == "" && req.Context.TenantID == "" {
		problem(w, r, http.StatusBadRequest, "TENANT_REQUIRED", "tenant_id is required", false)
		return
	}

	result, err := h.Service.Resolve(r.Context(), service.ResolutionRequest{
		TenantID:        req.TenantID,
		Context:         req.Context,
		Mappings:        req.Mappings,
		Bindings:        req.Bindings,
		EngineInstances: req.EngineInstances,
	})
	if err != nil {
		problem(w, r, http.StatusBadRequest, "RESOLUTION_FAILED", err.Error(), false)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"tenant_id": result.Context.TenantID,
		"mapping": map[string]any{
			"id":     result.Mapping.Mapping.ID,
			"status": result.Mapping.Mapping.Status,
		},
		"capability": map[string]any{
			"binding_mode":       result.Capability.BindingMode,
			"engine_instance_id": result.Capability.EngineInstanceID,
		},
		"policy": map[string]any{
			"allowed": result.Policy.Allowed,
			"reason":  result.Policy.Reason,
		},
		"topology": map[string]any{
			"id":          result.Topology.ID,
			"environment": result.Topology.Environment,
		},
	})
}
