package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
)

// CheckResponse is the response to a check request
type CheckResponse struct {
	Valid bool `json:"valid"`
}

// CheckOTPRequest is a request to check otp credentials for a user
type CheckOTPRequest struct {
	UserID string `json:"user_id"`
	OTP    string `json:"otp"`
}

// check otp checks if a given otp is valid
func checkOTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errors.WrongContentType(w)
		return
	}

	var req CheckOTPRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.MalformedBody(w)
		return
	}

	if req.UserID == "" {
		errors.MissingField(w, "user_id")
		return
	}

	if req.OTP == "" {
		errors.MissingField(w, "otp")
		return
	}

	user, err := database.GetUserByID(req.UserID)
	if err != nil {
		errors.UserDoesntExist(w)
		return
	}

	if !auth.IsOTPValid(user, req.OTP) {
		errors.InvalidCredentials(w)
		return
	}

	var resp = CheckResponse{
		Valid: true,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return

	}
}
