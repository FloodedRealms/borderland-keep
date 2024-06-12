package services

import (
	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
)

type UserService interface {
	//VerifyApiKey(providedKey string) (bool, error)
	CreateNewAPIUser(string, bool) (types.User, error)
	ValidateApiUser(providedClientId, providedApiKey string) (bool, error)
}

type UserServiceImpl struct {
	repo   repository.Repository
	logger util.Logger
}

func NewUserService(repo repository.Repository, logger util.Logger) *UserServiceImpl {
	return &UserServiceImpl{repo: repo, logger: logger}
}

func (us *UserServiceImpl) CreateNewAPIUser(friendlyName string, campaignLimit bool) (types.User, error) {
	newUser, err := types.GenerateNewUser(friendlyName)
	util.CheckErr(err)
	err = us.repo.SaveApiUser(newUser, campaignLimit)
	util.CheckErr(err)
	return newUser, nil
}
func (us *UserServiceImpl) ValidateApiUser(providedClientId, providedApiKey string) (bool, error) {
	userToValidate, err := us.repo.GetApiUserById(providedClientId, providedApiKey)
	if err != nil {
		return false, err
	}
	return userToValidate.Validate()
}
