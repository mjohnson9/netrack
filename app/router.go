package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

// This is the primary router used by all pages
var PrimaryRouter = mux.NewRouter()

func init() {
	http.Handle("/", PrimaryRouter)
}

func handleMethod(path string, handler http.Handler, method string) {
	PrimaryRouter.Path(path).Handler(handler).Methods(method)
}

func Get(path string, handler http.Handler) {
	handleMethod(path, handler, "GET")
}
