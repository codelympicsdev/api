package apiclients

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
	"github.com/gorilla/mux"
)

// GetAPIClientResponse is the response for getting a single api client
type GetAPIClientResponse struct {
	ID          string   `json:"id"`
	RedirectURL string   `json:"redirect_url"`
	Trusted     bool     `json:"trusted"`
	Scopes      []string `json:"scopes"`
}

// get a challenge
func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	client, err := database.GetAPIClientByID(vars["id"])
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			errors.NotFound(w)
			return
		}
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	var resp = &GetAPIClientResponse{
		ID:           client.ID,
		RedirectURL: client.RedirectURL,
		Trusted:      client.Trusted,
		Scopes:       client.Scopes,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}
}
