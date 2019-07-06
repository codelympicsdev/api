package auth

import "github.com/gorilla/mux"

// Route a auth request
func Route(r *mux.Router) {
	r.HandleFunc("/publickey", publickey).Methods("GET")
	r.HandleFunc("/signin", signin).Methods("POST")
	r.HandleFunc("/signup", signup).Methods("POST")
}
