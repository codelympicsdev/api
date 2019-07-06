package errors

import "net/http"

// Duplicate means there is a duplicate
func Duplicate(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "duplicate"}`))
}

// InternalServerError is when the server broke something
func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "internal_server_error"}`))
}
