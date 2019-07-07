package auth

import "github.com/gorilla/mux"

// AuthResponse is the response after signing in, up or upgrading
type AuthResponse struct {
	Token string `json:"token"`
}

// Route a auth request
func Route(r *mux.Router) {
	r.HandleFunc("/publickey", publickey).Methods("GET")

	r.HandleFunc("/signup", signup).Methods("POST")

	r.HandleFunc("/signin", signin).Methods("POST")
	r.HandleFunc("/upgrade/otp", upgradeWithOTP).Methods("POST")

	update := r.PathPrefix("/update").Subrouter()

	update.Use(TokenValidationMiddleware)
	update.HandleFunc("/password", updatePassword).Methods("POST")
	update.HandleFunc("/otp/enable", enableOTP).Methods("POST")
	update.HandleFunc("/otp/verify", verifyOTP).Methods("POST")

}
