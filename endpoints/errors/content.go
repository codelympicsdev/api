package errors

import "net/http"

// WrongContentType is the error for when the wrong content type was supplied
func WrongContentType(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "wrong_content_type"}`))
}

// MalformedBody means that the body supplied is not correct
func MalformedBody(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "malformed_body"}`))
}

// MissingField means that the content is missing a field
func MissingField(w http.ResponseWriter, field string) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "missing_field", "field": "` + field + `"}`))
}

// NotFound means that the content is not found
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "not_found"}`))
}

// MaxAttemptsReached means that the maximum attempts have been reached
func MaxAttemptsReached(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "max_attempts_reached"}`))
}

// AttemptTimedOut means that the attempt response was not recieved in time
func AttemptTimedOut(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"error": "attempt_timed_out"}`))
}
