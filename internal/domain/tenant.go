package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var resourceID = regexp.MustCompile(`^[a-z][a-z0-9]*(?:_[a-z0-9]+)*$`)
var regionID = regexp.MustCompile(`^[a-z]{2}-[a-z]+-[0-9]+$`)

type NotFoundError string

func (e NotFoundError) Error() string { return string(e) }

type RegisterTenant struct {
	LegalEntityID     string            `json:"legal_entity_id"`
	TenantID          string            `json:"tenant_id"`
	DisplayName       string            `json:"display_name"`
	IsolationStrategy string            `json:"isolation_strategy"`
	ResidencyRegion   string            `json:"residency_region"`
	RequestedProducts []string          `json:"requested_products,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

type ResolveContextRequest struct {
	TenantID  string `json:"tenant_id"`
	ProductID string `json:"product_id,omitempty"`
}

func (c ResolveContextRequest) Validate() error {
	if !validResource(c.TenantID) {
		return errors.New("tenant_id must be a canonical resource identifier")
	}
	if c.ProductID != "" && !validResource(c.ProductID) {
		return errors.New("product_id is invalid")
	}
	return nil
}

type ResolvedContext struct {
	TenantID           string `json:"tenant_id"`
	EntityID           string `json:"entity_id,omitempty"`
	LifecycleStatus    string `json:"lifecycle_status"`
	ProductID          string `json:"product_id,omitempty"`
	ProductEntitlement string `json:"product_entitlement,omitempty"`
}

type Tenant struct {
	TenantID          string            `json:"tenant_id"`
	LegalEntityID     string            `json:"legal_entity_id"`
	DisplayName       string            `json:"display_name"`
	IsolationStrategy string            `json:"isolation_strategy"`
	ResidencyRegion   string            `json:"residency_region"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	DesiredState      string            `json:"desired_state"`
	ObservedState     string            `json:"observed_state"`
	Revision          int64             `json:"revision"`
}

type EntitlementQuery struct {
	TenantID  string
	ProductID string
}

func (q EntitlementQuery) Validate() error {
	if !validResource(q.TenantID) {
		return errors.New("tenant_id must be a canonical resource identifier")
	}
	if !validResource(q.ProductID) {
		return errors.New("product_id must be a canonical resource identifier")
	}
	return nil
}

type Entitlement struct {
	TenantID  string `json:"tenant_id"`
	ProductID string `json:"product_id"`
	Status    string `json:"status"`
	Tier      string `json:"tier,omitempty"`
}

type LifecycleAction struct {
	TenantID string `json:"tenant_id"`
	Action   string `json:"action"`
}

func (a LifecycleAction) Validate() error {
	if !validResource(a.TenantID) {
		return errors.New("tenant_id must be a canonical resource identifier")
	}
	switch a.Action {
	case "activate", "suspend", "decommission":
		return nil
	default:
		return errors.New("action must be one of activate, suspend, or decommission")
	}
}

type LifecycleStatus string

const (
	LifecycleProvisioning  LifecycleStatus = "provisioning"
	LifecycleActive        LifecycleStatus = "active"
	LifecycleSuspended     LifecycleStatus = "suspended"
	LifecycleDecommissioning LifecycleStatus = "decommissioning"
	LifecycleDecommissioned  LifecycleStatus = "decommissioned"
)

func (s LifecycleStatus) Valid() bool {
	switch s {
	case LifecycleProvisioning, LifecycleActive, LifecycleSuspended, LifecycleDecommissioning, LifecycleDecommissioned:
		return true
	default:
		return false
	}
}

func TransitionLifecycle(from, to LifecycleStatus) (LifecycleStatus, bool) {
	if !from.Valid() || !to.Valid() {
		return "", false
	}
	if from == LifecycleDecommissioned && to != LifecycleDecommissioned {
		return "", false
	}
	allowed := map[LifecycleStatus][]LifecycleStatus{
		LifecycleProvisioning: {LifecycleActive, LifecycleSuspended},
		LifecycleActive:       {LifecycleSuspended, LifecycleDecommissioning},
		LifecycleSuspended:    {LifecycleActive, LifecycleDecommissioning},
		LifecycleDecommissioning: {LifecycleDecommissioned},
		LifecycleDecommissioned:  {LifecycleDecommissioned},
	}
	for _, next := range allowed[from] {
		if next == to {
			return to, true
		}
	}
	return "", false
}

func (c RegisterTenant) Validate() error {
	if !validResource(c.LegalEntityID) || !validResource(c.TenantID) {
		return errors.New("legal_entity_id and tenant_id must be canonical resource identifiers")
	}
	if n := len(strings.TrimSpace(c.DisplayName)); n < 1 || n > 255 {
		return errors.New("display_name must contain 1 to 255 characters")
	}
	if c.IsolationStrategy != "schema_per_tenant" && c.IsolationStrategy != "row_level_security" {
		return errors.New("isolation_strategy is invalid")
	}
	if !regionID.MatchString(c.ResidencyRegion) {
		return errors.New("residency_region is invalid")
	}
	if len(c.Metadata) > 20 {
		return errors.New("metadata may contain at most 20 entries")
	}
	seen := map[string]struct{}{}
	for _, product := range c.RequestedProducts {
		if !validResource(product) {
			return errors.New("requested_products contains an invalid identifier")
		}
		if _, ok := seen[product]; ok {
			return errors.New("requested_products must be unique")
		}
		seen[product] = struct{}{}
	}
	for _, value := range c.Metadata {
		if len(value) > 255 {
			return errors.New("metadata values may contain at most 255 characters")
		}
	}
	return nil
}
func validResource(v string) bool { return ValidResource(v) }

func ValidResource(v string) bool { return len(v) >= 3 && len(v) <= 63 && resourceID.MatchString(v) }

type Operation struct {
	OperationID string    `json:"operation_id"`
	TenantID    string    `json:"tenant_id"`
	State       string    `json:"state"`
	Revision    int       `json:"revision"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
