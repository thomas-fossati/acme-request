package models

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

type Delegation struct {
	ID             string     `json:"id"`
	CertURL        *string    `json:"cert-url,omitempty"`        // may be nil
	CreationDate   *time.Time `json:"creation-date,omitempty"`   // may be nil
	LastUpdate     *time.Time `json:"last-update,omitempty"`     // may be nil
	ExpirationDate *time.Time `json:"expiration-date,omitempty"` // may be nil
	CompletionDate *time.Time `json:"completion-date,omitempty"` // may be nil
	Status         string     `json:"status"`
	Details        *string    `json:"details,omitempty"` // may be nil
	// The following are the only ones allowed on input during creation.
	// Note that they are sanitized: CSR must be present, duration and certificate-lifetime
	// are checked against their configurable bounds.
	// Any other field is (silently) dropped (see also stores.sanitizeInDelegation())
	CSR          *string        `json:"csr"`
	CertLifetime *time.Duration `json:"certificate-lifetime"`
	Duration     *time.Duration `json:"duration"`
}

// ComputeETag computes a plausible ETag for the given Delegation resource.
// NOTE: this assumes the resource has one only representation (JSON) and that
// the JSON encoder is deterministic - which should be fine because the fields
// in a struct are static.
func ComputeETag(body []byte) string {
	h := sha1.New()
	h.Write(body)
	checkSum := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("\"%s\"", checkSum)
}
