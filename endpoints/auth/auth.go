package auth

import "github.com/gorilla/mux"

// AuthResponse is the response after signing in, up or upgrading
type AuthResponse struct {
	Token string `json:"token"`
}

// Route a auth request
func Route(r *mux.Router) {
	r.HandleFunc("/publickey", publickey).Methods("GET")
	r.HandleFunc("/signin", signin).Methods("POST")
	r.HandleFunc("/signup", signup).Methods("POST")
}
