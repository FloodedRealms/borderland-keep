package services

import (
	"fmt"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/floodedrealms/borderland-keep/types"
)

type UserService struct {
	repo   repository.Repository
	logger util.Logger
}

const userTable = "users"

func NewUserService(repo repository.Repository, logger util.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (us *UserService) CreateNewAPIUser(friendlyName string, campaignLimit bool) (types.User, error) {
	newUser, err := types.GenerateNewUser(friendlyName)
	util.CheckErr(err)
	err = us.repo.SaveApiUser(newUser, campaignLimit)
	util.CheckErr(err)
	return newUser, nil
}

func (us *UserService) CreateNewWebUser(username, email, password string) error {
	newUser, err := types.GenerateNewPasswordUser(username, password)
	util.CheckErr(err)
	saveStatment := fmt.Sprintf("INSERT INTO %s(name, email, campaigns_limited, password, salt) values(?, ?, ?, ?, ?)", userTable)
	_, err = us.repo.ExecuteQuery(saveStatment, newUser.Friendly_name, email, 1, newUser.RetreiveHash(), newUser.RetreiveSalt())
	util.CheckErr(err)
	return nil
}

func (us *UserService) ValidateApiUser(providedClientId, providedApiKey string) (bool, error) {
	userToValidate, err := us.repo.GetApiUserById(providedClientId, providedApiKey)
	if err != nil {
		return false, err
	}
	return userToValidate.Validate()
}

func (us *UserService) LimitUserCampaigns(id int) error {
	stmt := fmt.Sprintf("UPDATE %s SET campaigns_limited=1 WHERE id=%d", userTable, id)
	_, err := us.repo.ExecuteQuery(stmt)
	return err
}

func (us *UserService) UnlimitUserCampaigns(id int) error {
	stmt := fmt.Sprintf("UPDATE %s SET campaigns_limited=0 WHERE id=%d", userTable, id)
	_, err := us.repo.ExecuteQuery(stmt)
	return err
}
