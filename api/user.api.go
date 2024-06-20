package api

import (
	"errors"
	"net/http"

	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/types"
)

type UserApi struct {
	userService services.UserService
}

func NewUserApi(userService services.UserService) *UserApi {
	return &UserApi{userService: userService}
}

func (ua UserApi) ValidateApiUser(w http.ResponseWriter, r *http.Request) {
	var incomingUser *types.APIRequest
	err := decodeJSONBody(w, r, incomingUser)
	if err != nil {
		return
	}

	isValid, err := ua.userService.ValidateApiUser(incomingUser.Auth.ProvidedClientId, incomingUser.Auth.ProvidedAPIKey)
	if !isValid {
		http.Error(w, errors.New("user not valid").Error(), http.StatusForbidden)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (ua UserApi) RequireUserValidation(next http.Handler) http.Handler {
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
