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

// SignupRequest is what is used to sign up
type SignupRequest struct {
	ClientID string `json:"client_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// signup a user
func signup(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errors.WrongContentType(w)
		return
	}

	var req SignupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.MalformedBody(w)
		return
	}

	if req.ClientID == "" {
		errors.MissingField(w, "client_id")
		return
	}

	if req.FullName == "" {
		errors.MissingField(w, "full_name")
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

	user, err := auth.Signup(req.FullName, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			errors.Duplicate(w)
			return
		}
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	token := auth.NewToken(user, client)
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
