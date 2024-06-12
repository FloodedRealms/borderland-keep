package commands

import (
	"log"

	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/util"
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
