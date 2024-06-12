package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/floodedrealms/adventure-archivist/api"
	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/util"
)

func main() {
	//	memRepo, err := repository.NewMemoryRepo()
	//	util.CheckErr(err)

	logger := util.NewLogger(true)

	sqlRepo, err := repository.NewSqliteRepo("archivist.db", logger)
	util.CheckErr(err)

	campaignService := services.NewCampaignService(sqlRepo, logger, context.TODO())
	characterService := services.NewCharacterService(sqlRepo, logger, context.TODO())
	adventureRecordService := services.NewAdventureRecordService(sqlRepo, context.TODO())

	campaignApi := api.NewCampaignApi(campaignService, characterService)
	adventureRecordApi := api.NewAdventureRecordApi(adventureRecordService)
	characterApi := api.NewCharacterApi(characterService)

	//characterService := services.NewCharacterService(memRepo, context.TODO())
	//characterApi := api.NewCharacterApi(characterService)

	router := gin.Default()

	//Campaign Endpoints
	router.POST("/api/campaigns", campaignApi.CreateCampaign)
	router.POST("/api/campaigns/:campaignId/adventures", adventureRecordApi.CreateAdventureRecord)
	router.POST("/api/campaigns/:campaignId/characters", characterApi.CreateCharacterForCampaign)

	router.PATCH("/api/campaigns/:campaignId", campaignApi.UpdateCampaign)

	router.GET("/api/campaigns", campaignApi.ListCampaigns)
	router.GET("/api/campaigns/:campaignId", campaignApi.GetCampaign)
	router.GET("/api/campaigns/:campaignId/adventures", adventureRecordApi.ListAdventureRecordsForCampaign)

	router.DELETE("/api/campaigns/:campaignId", campaignApi.DeleteCampaign)

	router.OPTIONS("/api/campaigns", preflight)
	router.OPTIONS("/api/campaigns/:campaignId", preflight)

	//Adventure Endpoints
	//router.POST("/adventures/:adventureId/loot", adventureRecordApi.AddLootToAdventure)
	router.POST("/api/adventures/:adventureId/characters/:characterId", characterApi.ManageCharactersForAdventure)

	router.PATCH("/api/adventures/:adventureId", adventureRecordApi.UpdateAdventure)

	router.GET("/api/adventures/:adventureId", adventureRecordApi.GetAdventure)
	// router.GET("/adventures/:adventureId/experience", adventureRecordApi.GetAdventureExperience)

	//	router.DELETE("/adventures/{adventureId}", adventureRecordApi.DeleteAdventure)

	router.OPTIONS("/api/adventures", preflight)
	router.OPTIONS("/api/adventures/:adventureId", preflight)
	//router.OPTIONS("/adventures/:adventureId/loot/:type", preflight)
	router.OPTIONS("/api/adventures/:adventureId/combat", preflight)
	router.OPTIONS("/api/adventures/:adventureId/:characterId/:op", preflight)

	//Character Endpoints
	router.GET("/api/characters/:characterId", characterApi.GetCharacterById)

	//router.PATCH("/characters/:characterId", characterApi.UpdateCharacter)

	//	router.DELETE("/characters/:characterId", characterApi.DeleteCharacter)

	router.OPTIONS("/api/characters", preflight)
	router.OPTIONS("/api/characters/:adventureId", preflight)

	router.Run("localhost:9090")
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
	c.JSON(http.StatusOK, struct{}{})
}
