package apperrors

import (
	"encoding/json"
	"net/http"

	"github.com/moogar0880/problems"
)

// ReportProblem takes an error and maps it to an appropriate RFC7807 "Problem Detail"
// (The most simplistic approach for now.)
func ReportProblem(err error) (int, string, []byte) {
	var code int

	switch err {
	case ErrDelegationUnknown:
		code = http.StatusNotFound
	case ErrDelegationMissingParameter, ErrDelegationBadParameter:
		code = http.StatusBadRequest
	// TODO(tho) Other mappings
	default:
		code = http.StatusInternalServerError
	}

	var body []byte

	body, jerr := json.Marshal(problems.NewStatusProblem(code))
	if jerr != nil {
		return http.StatusInternalServerError, "text/plain", []byte("JSON encoding failed")
	}

	return code, problems.ProblemMediaType, body
}
