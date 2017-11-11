package apis

import (
	"log"
	"net/http"

	"../apperrors"
	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
)

func Start(listenSpec string, router *mux.Router) {
	go http.ListenAndServe(listenSpec, router)
}

// RequestHandler is the signature of all API endpoint handler functions
type RequestHandler func(http.ResponseWriter, *http.Request) error

// A RequestHandler implements ServeHTTP and is therefore a http.Handler
// This wrapper aggregates all error handling at the API layer.
func (fn RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		status, ct, body := apperrors.ReportProblem(e)

		w.Header().Set("Content-Type", ct)
		w.WriteHeader(status)
		w.Write(body)
	}
}

// RequestLogger uses httpsnoop to collect metrics on our request handler and
// print an access.log-like line per transaction
func RequestLogger(h RequestHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		m := httpsnoop.CaptureMetrics(h, w, req)
		log.Printf(
			"%s\t%s\t%d\t%s\t%dbytes",
			req.Method,
			req.URL,
			m.Code,
			m.Duration,
			m.Written,
		)
	})
}
