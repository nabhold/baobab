package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nabhold/baobab-cp/internal/domain"
	"github.com/nabhold/baobab-cp/internal/resolver"
)

// ResolutionRequest is the service-level input for a control-plane resolution request.
type ResolutionRequest struct {
	TenantID        string
	Context         resolver.Context
	Mappings        []domain.Mapping
	Bindings        []resolver.CapabilityBinding
	EngineInstances []resolver.EngineInstance
}

// ResolutionResult contains the fully resolved runtime decision for a tenant request.
type ResolutionResult struct {
	Context    resolver.Context
	Mapping    resolver.ResolvedMapping
	Capability resolver.ResolvedCapability
	Policy     resolver.PolicyDecision
	Topology   resolver.EngineInstance
}

// ResolutionService exposes the composed resolver pipeline as a service interface.
type ResolutionService struct {
	Pipeline resolver.ResolutionPipeline
}

func (s ResolutionService) Resolve(ctx context.Context, req ResolutionRequest) (ResolutionResult, error) {
	if ctx == nil {
		return ResolutionResult{}, errors.New("context is required")
	}
	if req.TenantID == "" && req.Context.TenantID == "" {
		return ResolutionResult{}, errors.New("tenant_id is required")
	}
	if req.Context.TenantID == "" {
		req.Context.TenantID = req.TenantID
	}

	pipelineResult, err := s.Pipeline.Resolve(ctx, resolver.ResolutionRequest{
		TenantID:        req.TenantID,
		Context:         req.Context,
		Candidates:      req.Mappings,
		Bindings:        req.Bindings,
		EngineInstances: req.EngineInstances,
	})
	if err != nil {
		return ResolutionResult{}, fmt.Errorf("resolution failed: %w", err)
	}

	return ResolutionResult{
		Context:    pipelineResult.Context,
		Mapping:    pipelineResult.Mapping,
		Capability: pipelineResult.Capability,
		Policy:     pipelineResult.Policy,
		Topology:   pipelineResult.Topology,
	}, nil
}
