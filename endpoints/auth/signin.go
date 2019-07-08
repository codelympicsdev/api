package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
)

// SigninRequest is what is used to sign in
type SigninRequest struct {
	ClientID string `json:"client_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// signin a user
func signin(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errors.WrongContentType(w)
		return
	}

	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.MalformedBody(w)
		return
	}

	if req.ClientID == "" {
		errors.MissingField(w, "client_id")
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

	client, err := database.GetAPIClientByID(req.ClientID)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	user, err := auth.Signin(req.Email, req.Password)
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "invalid username or password") {
			errors.InvalidCredentials(w)
			return
		}
		errors.InternalServerError(w)
		return
	}

	var token *auth.Token
	if user.OTPEnabled {
		token = auth.NewUnverifiedToken(user, client)
	} else {
		token = auth.NewToken(user, client)
	}

	t, err := token.Sign()
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	var resp = AuthResponse{
		Token: t,
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}
}
