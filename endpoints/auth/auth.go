package auth

import "github.com/gorilla/mux"

// Route a auth request
func Route(r *mux.Router) {
	r.HandleFunc("/publickey", publickey).Methods("GET")

	protected := r.PathPrefix("").Subrouter()
	protected.Use(RootTrustClientValidationMiddleware)
	protected.HandleFunc("/signin/emailpassword", signinEmailPassword).Methods("POST")
	protected.HandleFunc("/check/otp", checkOTP).Methods("POST")
	protected.HandleFunc("/generatetoken", generateToken).Methods("POST")
	// protected.HandleFunc("/signup", signup).Methods("POST")

	update := r.PathPrefix("/update").Subrouter()
	update.Use(TokenValidationMiddleware)
	update.HandleFunc("/password", updatePassword).Methods("POST")
	update.HandleFunc("/otp/enable", enableOTP).Methods("POST")
	update.HandleFunc("/otp/verify", verifyOTP).Methods("POST")

}
