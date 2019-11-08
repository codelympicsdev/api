package apiclients

import (
	"github.com/codelympicsdev/api/endpoints/auth"
	"github.com/gorilla/mux"
)

// Route a auth request
func Route(r *mux.Router) {
	protected := r.PathPrefix("").Subrouter()
	protected.Use(auth.RootTrustClientValidationMiddleware)
	protected.HandleFunc("/{id}", get).Methods("GET")
}
