package apperrors

import (
	"errors"
)

var (
	// ErrDelegationUnknown is returned if the supplied Delegation is not known to the store layer
	ErrDelegationUnknown = errors.New("Delegation unknown")
	// ErrDelegationMissingParameter is returned if the supplied Delegation is missing a non-optional parameter
	ErrDelegationMissingParameter = errors.New("Missing non optional parameter")
	// ErrDelegationBadParameter is returned if the syntax or semantics of one or more input parameters of the Delegation request is wrong
	ErrDelegationBadParameter = errors.New("Bad parameter supplied")
)
