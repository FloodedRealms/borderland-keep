package api

import (
	"net/http"

	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
	userService services.UserService
}

func NewUserApi(userService services.UserService) *UserApi {
	return &UserApi{userService: userService}
}

func (ua *UserApi) ValidateApiUser(ctx *gin.Context) {
	var incomingUser *types.APIRequest
	if err := ctx.ShouldBindJSON(&incomingUser); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	isValid, err := ua.userService.ValidateApiUser(incomingUser.Auth.ProvidedClientId, incomingUser.Auth.ProvidedAPIKey)
	if !isValid {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "access-denied", "message": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "valid", "message": "User is valid"})
}

func (ua *UserApi) RequireUserValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		givenClientId := ctx.Request.Header.Get("X-Archivist-Client-Id")
		givenAPIKey := ctx.Request.Header.Get("X-Archivist-API-Key")

		isValid, err := ua.userService.ValidateApiUser(givenClientId, givenAPIKey)
		if !isValid {
			ctx.JSON(http.StatusForbidden, gin.H{"status": "access-denied", "message": err})
			return
		}
		ctx.Next()
	}

}
