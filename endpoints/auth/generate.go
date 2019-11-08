package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
)

// GenerateTokenResponse is the response for generating a token
type GenerateTokenResponse struct {
	Token string `json:"token"`
}

// GenerateTokenRequest is a request to generate a token
type GenerateTokenRequest struct {
	ClientID string   `json:"client_id"`
	UserID   string   `json:"user_id"`
	Scopes   []string `json:"scopes"`
}

// generate token generates a token
func generateToken(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errors.WrongContentType(w)
		return
	}

	var req GenerateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.MalformedBody(w)
		return
	}

	if req.ClientID == "" {
		errors.MissingField(w, "client_id")
		return
	}

	if req.UserID == "" {
		errors.MissingField(w, "user_id")
		return
	}

	if req.Scopes == nil || len(req.Scopes) == 0 {
		errors.MissingField(w, "scopes")
		return
	}

	user, err := database.GetUserByID(req.UserID)
	if err != nil {
		errors.UserDoesntExist(w)
		return
	}

	client, err := database.GetAPIClientByID(req.ClientID)
	if err != nil {
		errors.APIClientDoesntExist(w)
		return
	}

	token := auth.NewToken(user, client, req.Scopes)
	t, err := token.Sign()
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	var resp = GenerateTokenResponse{
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
