package services

import (
	"fmt"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/util"
)

type UserService struct {
	repo   repository.Repository
	logger util.Logger
}

const userTable = "users"

func NewUserService(repo repository.Repository, logger util.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

/*
	func (us *UserService) CreateNewAPIUser(friendlyName string, campaignLimit bool) (types.User, error) {
		newUser, err := types.GenerateNewUser(friendlyName)
		util.CheckErr(err)
		err = us.repo.SaveApiUser(newUser, campaignLimit)
		util.CheckErr(err)
		return newUser, nil
	}
*/

func (us *UserService) IsNameTaken(name string) (bool, error) {
	selectStatment := fmt.Sprintf("SELECT count(id) FROM %s WHERE name=?;", userTable)
	rows, err := us.repo.RunQuery(selectStatment, name)
	if err != nil {
		return true, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		rows.Scan(&count)
	}
	return count > 0, nil
}

func (us *UserService) InsertWebUser(username, email, passwordHash, salt string) error {
	saveStatment := fmt.Sprintf("INSERT INTO %s(name, email, campaigns_limited, password, salt) values(?, ?, ?, ?, ?)", userTable)
	_, err := us.repo.ExecuteQuery(saveStatment, username, email, 1, passwordHash, salt)
	util.CheckErr(err)
	return nil
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
