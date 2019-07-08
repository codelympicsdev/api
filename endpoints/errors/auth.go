package errors

import "net/http"

// InvalidCredentials means that the supplied credentials are not correct
func InvalidCredentials(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "invalid_credentials"}`))
}

// TokenAlreadyUpgraded means that the supplied token is already upgraded
func TokenAlreadyUpgraded(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "token_already_upgraded"}`))
}

// UnverifiedToken means that the supplied token is not verfied yet
func UnverifiedToken(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "unverified_token"}`))
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

// InvalidAPIClient means that the supplied api client (credentials) are not correct
func InvalidAPIClient(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "invalid_api_client"}`))
}
