package router

import (
	"github.com/gorilla/mux"
)

// NewRouter creates a new request router instance, which will be passed as a
// parameter to the server.  The request router receives all HTTP connections
// and passes them on to the request handlers you will register on it.
func NewRouter() *mux.Router {
	return mux.NewRouter().StrictSlash(true)
}
