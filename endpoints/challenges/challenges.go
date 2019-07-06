package challenges

import (
	"github.com/codelympicsdev/api/endpoints/auth"
	"github.com/gorilla/mux"
)

// Route a auth request
func Route(r *mux.Router) {
	r.Use(auth.TokenValidationMiddleware)
	r.HandleFunc("/{id}", get).Methods("GET")
}
