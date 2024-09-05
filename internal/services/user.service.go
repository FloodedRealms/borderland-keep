package services

import (
	"fmt"
	"time"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/util"
)

type UserService struct {
	repo   repository.Repository
	logger util.Logger
}

const userTable = "users"
const sessionTable = "sessions"

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

type UserNotFoundError struct {
	username string
}

func (u UserNotFoundError) Error() string {
	return fmt.Sprintf("Username %s was not found", u.username)
}

func (us UserService) RetrieveWebUserInformation(username string) (id string, password, salt string, err error) {
	selectStatment := fmt.Sprintf("SELECT  id, password, salt FROM %s WHERE name=?;", userTable)
	rows, err := us.repo.RunQuery(selectStatment, username)
	if err != nil {
		return "", "", "", err
	}
	defer rows.Close()

	usercount := 0
	for rows.Next() {
		usercount++
		rows.Scan(&id, &password, &salt)
	}
	if usercount < 1 {
		return "", "", "", UserNotFoundError{username: username}
	}
	if usercount > 1 {
		return "", "", "", fmt.Errorf("too many results")
	}
	return id, password, salt, nil
}

func (us UserService) StoreSession(token, userId, userName string, expiryTime time.Time) error {
	deleteStatement := fmt.Sprintf("DELETE FROM %s WHERE user_id=?;", sessionTable)
	deleteParams := []interface{}{userId}
	saveStatment := fmt.Sprintf("INSERT INTO %s(uuid, user_id, user_name, expiry_time) values(?, ?, ?, ?);", sessionTable)
	saveParams := []interface{}{token, userId, userName, expiryTime}
	err := us.repo.DoTransaction([]string{deleteStatement, saveStatment}, [][]interface{}{deleteParams, saveParams})
	return err
}

func (us UserService) DeleteSession(token string) error {
	deleteStatement := fmt.Sprintf("DELETE FROM %s WHERE uuid=?;", sessionTable)
	_, err := us.repo.ExecuteQuery(deleteStatement, token)
	return err
}

func (us UserService) RetrieveSession(sessionId string) (userId, userName string, expiryTime time.Time, exists bool, err error) {
	selectStatment := fmt.Sprintf("SELECT user_id, user_name, expiry_time FROM %s WHERE uuid=?;", sessionTable)
	rows, err := us.repo.RunQuery(selectStatment, sessionId)
	if err != nil {
		return "", "", time.Now(), false, err
	}
	defer rows.Close()

	sessionCount := 0
	for rows.Next() {
		sessionCount++
		rows.Scan(&userId, &userName, &expiryTime)
	}
	if sessionCount < 1 {
		return "", "", time.Now(), false, nil
	}
	return userId, userName, expiryTime, true, nil
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
