package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// GetUserResponse is the response for a single user
type GetUserResponse struct {
	ID         string `json:"id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	AvatarURL  string `json:"avatar_url"`
	OTPEnabled *bool  `json:"otp_enabled,omitempty"`
}

// get a user
func get(w http.ResponseWriter, r *http.Request) {
	token, ok := context.Get(r, "token").(*auth.Token)
	if ok == false {
		errors.InternalServerError(w)
		return
	}

	vars := mux.Vars(r)

	if (vars["id"] == token.ID && (token.HasScope("user.basic") || token.HasScope("user.read"))) || token.HasScope("admin.users") {
		u, err := database.GetUserByID(vars["id"])
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}

		resp := &GetUserResponse{
			ID:         u.ID,
			FullName:   u.FullName,
			Email:      u.Email,
			AvatarURL:  u.AvatarURL,
			OTPEnabled: nil,
		}

		if vars["id"] == token.ID && token.HasScope("user.read") {
			resp.OTPEnabled = &u.OTPEnabled
		}

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}
	} else {
		if vars["id"] == token.ID {
			errors.MissingPermission(w, "user.basic")
		} else {
			errors.MissingPermission(w, "admin.users")
		}
		return
	}
}
