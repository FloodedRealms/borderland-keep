package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/google/uuid"
)

func TestIsNameTaken(t *testing.T) {
	username := "TestUsernameIsTakenDoNotPutThisInDataseForever"
	logger := util.NewLogger(true)
	sqlRepo, err := repository.NewSqliteRepo(logger, "test.db")
	if err != nil {
		t.Errorf("Unable to get test DB: %s", err.Error())
	}
	insertUser := "INSERT INTO users(name, email, campaigns_limited, password, salt) values(?, ?, 1, ?, ?);"
	deleteUser := "DELETE FROM users WHERE name=?;"
	_, err = sqlRepo.ExecuteQuery(insertUser, username, "test@test.xyz", "password", "salt")
	if err != nil {
		t.Errorf("Unable to insert test user DB: %s", err.Error())
		return
	}
	expected := true
	service := services.NewUserService(sqlRepo, *logger)
	taken, err := service.IsNameTaken(username)
	if err != nil {
		t.Errorf("Unable to get username from DB: %s", err.Error())
		return
	}
	if taken != expected {
		t.Errorf("Username should have been taken")

	}
	_, err = sqlRepo.ExecuteQuery(deleteUser, username)
	if err != nil {
		t.Errorf("Unable to delete test user DB: %s", err.Error())
		return
	}
	expected = false
	taken, err = service.IsNameTaken(username)
	if err != nil {
		t.Errorf("Unable to get username from DB: %s", err.Error())
		return
	}
	if taken != expected {
		t.Errorf("Username should NOT have been taken")
	}

}

func TestUserhasEditAccessToCampaign(t *testing.T) {
	username := "testAccessUser"
	userId := "1"
	uuid := uuid.New().String()
	accessibleCampaign := 1
	inaccessibleCampaign := 2
	logger := util.NewLogger(true)
	sqlRepo, err := repository.NewSqliteRepo(logger, "test.db")
	if err != nil {
		t.Errorf("Unable to get test DB: %s", err.Error())
	}
	service := services.NewUserService(sqlRepo, *logger)
	service.StoreSession(uuid, userId, username, time.Now().Add(120*time.Second))

	expected := true
	got, err := service.UserhasEditAccessToCampaign(uuid, accessibleCampaign)
	if err != nil {
		t.Errorf("Unable to get test DB: %s", err.Error())
	}
	if expected != got {
		t.Error(fmt.Sprintf("Expected userId %s to have access to campaignId %d!", userId, accessibleCampaign))
	}
	expected = false
	got, err = service.UserhasEditAccessToCampaign(uuid, inaccessibleCampaign)
	if err != nil {
		t.Errorf("Unable to get test DB: %s", err.Error())
	}
	if expected != got {
		t.Error(fmt.Sprintf("Expected userId %s to NOT have access to campaignId %d!", userId, inaccessibleCampaign))
	}
}
