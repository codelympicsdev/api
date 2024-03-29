package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/endpoints/errors"
)

// SigninResponse the response to a signin request
type SigninResponse struct {
	UserID      string `json:"user_id"`
	Requires2FA bool   `json:"requires_2fa"`
}

// SigninEmailPasswordRequest is a request that gets a user and checks their credentials by email and password
type SigninEmailPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// check a user with email and password
func signinEmailPassword(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errors.WrongContentType(w)
		return
	}

	var req SigninEmailPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.MalformedBody(w)
		return
	}

	if req.Email == "" {
		errors.MissingField(w, "email")
		return
	}

	if req.Password == "" {
		errors.MissingField(w, "password")
		return
	}

	user, err := auth.SigninEmailPassword(req.Email, req.Password)
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "invalid username or password") {
			errors.InvalidCredentials(w)
			return
		}
		errors.InternalServerError(w)
		return
	}

	var resp = SigninResponse{
		UserID:      user.ID,
		Requires2FA: user.OTPEnabled,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}
}
