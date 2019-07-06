package errors

import "net/http"

// InvalidCredentials means that the supplied credentials are not correct
func InvalidCredentials(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "invalid_credentials"}`))
}

// MissingPermission means that the token is missing a permission
func MissingPermission(w http.ResponseWriter, scope string) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "missing_permission", "scope": "` + scope + `"}`))
}
