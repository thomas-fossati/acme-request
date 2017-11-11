package models

import "time"

type Delegation struct {
	ID             string     `json:"id"`
	CertURL        *string    `json:"cert-url,omitempty"` // may be Null
	CreationDate   time.Time  `json:"creation-date"`
	LastUpdate     *time.Time `json:"last-update,omitempty"`     // may be Null
	ExpirationDate *time.Time `json:"expiration-date,omitempty"` // may be Null
	CompletionDate *time.Time `json:"completion-date,omitempty"` // may be Null
	Status         string     `json:"status"`
	Details        *string    `json:"details,omitempty"` // may be Null
	// The following are the only ones allowed on input during creation.
	// Note that they are sanitized: CSR must be present, duration and certificate-lifetime
	// are checked against their configurable bounds.
	// Any other field is (silently) dropped (see also stores.sanitizeInDelegation())
	CSR          *string        `json:"csr"`
	CertLifetime *time.Duration `json:"certificate-lifetime"`
	Duration     *time.Duration `json:"duration"`
}

// Compute ETag
func ComputeETag(d *Delegation) string {
	// TODO(tho)
	// - Create a consistent linearisation of the fields that are subject to change
	// - hash it
	// - truncate as needed
	return "\"xyzzy\""
}
