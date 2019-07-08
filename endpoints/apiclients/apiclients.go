package apiclients

import (
	"github.com/gorilla/mux"
)

// Route a auth request
func Route(r *mux.Router) {
	r.HandleFunc("/{id}", get).Methods("GET")
}
