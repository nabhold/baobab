package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var resourceID = regexp.MustCompile(`^[a-z][a-z0-9]*(?:_[a-z0-9]+)*$`)
var regionID = regexp.MustCompile(`^[a-z]{2}-[a-z]+-[0-9]+$`)

type RegisterTenant struct {
	LegalEntityID     string            `json:"legal_entity_id"`
	TenantID          string            `json:"tenant_id"`
	DisplayName       string            `json:"display_name"`
	IsolationStrategy string            `json:"isolation_strategy"`
	ResidencyRegion   string            `json:"residency_region"`
	RequestedProducts []string          `json:"requested_products,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
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
func validResource(v string) bool { return len(v) >= 3 && len(v) <= 63 && resourceID.MatchString(v) }

type Operation struct {
	OperationID string    `json:"operation_id"`
	TenantID    string    `json:"tenant_id"`
	State       string    `json:"state"`
	Revision    int       `json:"revision"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
