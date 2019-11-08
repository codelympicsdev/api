package errors

import "net/http"

// InvalidCredentials means that the supplied credentials are not correct
func InvalidCredentials(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "invalid_credentials"}`))
}

// UserDoesntExist means that the user doesnt exist
func UserDoesntExist(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "user_doesnt_exist"}`))
}

// OTPAlreadyEnabled means otp is already enabled
func OTPAlreadyEnabled(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "otp_already_enabled"}`))
}

// MissingPermission means that the token is missing a permission
func MissingPermission(w http.ResponseWriter, scope string) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "missing_permission", "scope": "` + scope + `"}`))
}

// APIClientDoesntExist means that the api client doesnt exist
func APIClientDoesntExist(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "api_client_doesnt_exist"}`))
}

// InvalidRootTrustClient means that the supplied root trust client (credentials) are not correct
func InvalidRootTrustClient(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "invalid_root_trust_client"}`))
}
