package guardsman

import (
	"log"
	"strconv"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
)

func createApiUser(friendlyName string) {
	logger := util.NewLogger(true)
	sqlRepo, err := repository.NewSqliteRepo(logger)
	util.CheckErr(err)
	userService := services.NewUserService(sqlRepo, *logger)
	newUser, err := userService.CreateNewAPIUser(friendlyName, true)
	util.CheckErr(err)
	log.Printf("New User created with Id: %s api_key: %s. Friendly name is: %s", newUser.DisplayUUID(), newUser.DisplayAPIKey(), newUser.DisplayUserName())
}

func createWebUser(username, email, password string) {
	logger := util.NewLogger(true)
	sqlRepo, err := repository.NewSqliteRepo(logger)
	util.CheckErr(err)
	userService := services.NewUserService(sqlRepo, *logger)
	err = userService.CreateNewWebUser(username, email, password)
	util.CheckErr(err)
	log.Printf("New User created with. Friendly name is: %s", username)
}

func CreateUser(userType, friendlyName, email, password string) {
	switch userType {
	case "api":
		createApiUser(friendlyName)
	case "web":
		createWebUser(friendlyName, email, password)
	default:
		createWebUser(friendlyName, email, password)
	}
}

func UnlimitUserCampaigns(id string) error {
	userId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	logger := util.NewLogger(true)
	sqlRepo, err := repository.NewSqliteRepo(logger)
	util.CheckErr(err)
	userService := services.NewUserService(sqlRepo, *logger)
	return userService.UnlimitUserCampaigns(userId)
}

func LimitUserCampaigns(id string) error {
	userId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	logger := util.NewLogger(true)
	sqlRepo, err := repository.NewSqliteRepo(logger)
	util.CheckErr(err)
	userService := services.NewUserService(sqlRepo, *logger)
	return userService.LimitUserCampaigns(userId)
}
