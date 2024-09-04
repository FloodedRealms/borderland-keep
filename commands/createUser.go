package commands

import (
	"log"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
)

func createApiUser(friendlyName string, campaignLimit bool) {
	logger := util.NewLogger(true)
	sqlRepo, err := repository.NewSqliteRepo("archivist.db", logger)
	util.CheckErr(err)
	userService := services.NewUserService(sqlRepo, *logger)
	newUser, err := userService.CreateNewAPIUser(friendlyName, campaignLimit)
	util.CheckErr(err)
	log.Printf("New User created with Id: %s api_key: %s. Friendly name is: %s", newUser.DisplayUUID(), newUser.DisplayAPIKey(), newUser.DisplayUserName())
}

func CreateUser(userType, friendlyName string, campaignLimit bool) {
	switch userType {
	case "api":
		createApiUser(friendlyName, campaignLimit)
	}
}
