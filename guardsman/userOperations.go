package guardsman

import (
	"fmt"
	"log"
	"strconv"

	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
)

/*
	func createApiUser(friendlyName string) {
		logger := util.NewLogger(true)
		sqlRepo, err := repository.NewSqliteRepo(logger)
		util.CheckErr(err)
		userService := services.NewUserService(sqlRepo, *logger)
		newUser, err := userService.CreateNewAPIUser(friendlyName, true)
		util.CheckErr(err)
		log.Printf("New User created with Id: %s api_key: %s. Friendly name is: %s", newUser.DisplayUUID(), newUser.DisplayAPIKey(), newUser.DisplayUserName())
	}
*/

func (g Guardsman) createWebUser(username, email, password string) {
	nameTaken, err := g.userService.IsNameTaken(username)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if nameTaken {
		fmt.Println("Username must be Unique")
		return
	}
	user, err := GenerateNewPasswordUser(username, password)
	util.CheckErr(err)
	user.Email = email
	err = g.userService.InsertWebUser(user.Friendly_name, user.Email, user.RetreiveHash(), user.RetreiveSalt())
	util.CheckErr(err)
	log.Printf("New User created with. Friendly name is: %s", username)
}

func (g Guardsman) CreateUser(userType, friendlyName, email, password string) {
	switch userType {
	case "api":
		//createApiUser(friendlyName)
	case "web":
		g.createWebUser(friendlyName, email, password)
	default:
		g.createWebUser(friendlyName, email, password)

	}
}

func (g Guardsman) UnlimitUserCampaigns(id string) error {
	userId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	logger := util.NewLogger(true)
	userService := services.NewUserService(g.repo, *logger)
	return userService.UnlimitUserCampaigns(userId)
}

func (g Guardsman) LimitUserCampaigns(id string) error {
	userId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	logger := util.NewLogger(true)
	userService := services.NewUserService(g.repo, *logger)
	return userService.LimitUserCampaigns(userId)
}
