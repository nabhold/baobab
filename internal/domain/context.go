package domain

import (
	"errors"
	"regexp"
	"time"
)

var productID = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[-_][a-z0-9]+)*$`)

type ResolveContext struct {
	ProductID string `json:"product_id"`
}

func (c ResolveContext) Validate() error {
	if len(c.ProductID) < 3 || len(c.ProductID) > 63 || !productID.MatchString(c.ProductID) {
		return errors.New("product_id must be a canonical product identifier")
	}
	return nil
}

type ResolvedContext struct {
	TenantID        string    `json:"tenant_id"`
	EntityID        string    `json:"entity_id"`
	LifecycleStatus string    `json:"lifecycle_status"`
	ProductID       string    `json:"product_id"`
	Entitled        bool      `json:"entitled"`
	EntitlementTier *string   `json:"entitlement_tier,omitempty"`
	CacheTTLSeconds int       `json:"cache_ttl_seconds"`
	ResolvedAt      time.Time `json:"resolved_at"`
	CorrelationID   string    `json:"correlation_id"`
}
