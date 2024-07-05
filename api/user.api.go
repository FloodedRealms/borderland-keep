package api

import (
	"errors"
	"net/http"

	"github.com/floodedrealms/adventure-archivist/internal/services"
	"github.com/floodedrealms/adventure-archivist/types"
)

type User struct {
	userService services.UserService
}

func (ua User) ValidatePassword(w http.ResponseWriter, r *http.Request) {
	var incomingReques *types.APIRequest
	err := decodeJSONBody(w, r, incomingReques)
	if err != nil {
		return
	}

	isValid, err := ua.userService.ValidateApiUser(incomingReques.Auth.ProvidedClientId, incomingReques.Auth.ProvidedAPIKey)
	if !isValid {
		http.Error(w, errors.New("client not valid").Error(), http.StatusForbidden)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (ua User) RequirePassword(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		givenClientId := r.Header[http.CanonicalHeaderKey("X-Archivist-Client-Id")][0]
		givenAPIKey := r.Header[http.CanonicalHeaderKey("X-Archivist-API-Key")][0]

		isValid, err := ua.userService.ValidateApiUser(givenClientId, givenAPIKey)
		if !isValid {
			http.Error(w, errors.New("user not valid").Error(), http.StatusForbidden)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		next.ServeHTTP(w, r)
	})

}
