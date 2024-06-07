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
	campaignApi := api.NewCampaignApi(campaignService)

	adventureRecordService := services.NewAdventureRecordService(sqlRepo, context.TODO())
	adventureRecordApi := api.NewAdventureRecordApi(adventureRecordService)

	//characterService := services.NewCharacterService(memRepo, context.TODO())
	//characterApi := api.NewCharacterApi(characterService)

	router := gin.Default()

	//Campaign Endpoints
	router.POST("/campaigns", campaignApi.CreateCampaign)
	router.POST("/campaigns/:campaignId/adventures", adventureRecordApi.CreateAdventureRecord)

	router.GET("/campaigns", campaignApi.ListCampaigns)
	router.GET("/campaigns/:campaignId", campaignApi.GetCampaign)
	router.GET("/campaigns/:campaignId/adventures", adventureRecordApi.ListAdventureRecordsForCampaign)

	router.DELETE("/campaigns/:campaignId", campaignApi.DeleteCampaign)

	router.OPTIONS("/campaigns", preflight)
	router.OPTIONS("/campaigns/:campaignId", preflight)

	//Adventure Endpoints

	router.POST("/adventures/:adventureId/loot", /*func(c *gin.Context) {
			id := c.Param("adventureIdid")
			c.String(http.StatusOK, id)
			// Handle request with id parameter
		}) */adventureRecordApi.AddLootToAdventure)
	router.GET("/adventures/:adventureId", adventureRecordApi.GetAdventure)
	// router.GET("/adventures/:adventureId/experience", adventureRecordApi.GetAdventureExperience)

	// router.POST("/adventures/:adventureId/combat", adventureRecordApi.AddCombatToAdventure)
	// router.POST("/adventures/:adventureId/:characterId/add", adventureRecordApi.AddCharacterToAdventure)
	// router.POST("/adventures/:adventureId/:characterId/remove", adventureRecordApi.RemoveCharacterFromAdventure)

	//	router.DELETE("/adventures/{adventureId}", adventureRecordApi.DeleteAdventure)

	router.OPTIONS("/adventures", preflight)
	router.OPTIONS("/adventures/:adventureId", preflight)
	//router.OPTIONS("/adventures/:adventureId/loot/:type", preflight)
	router.OPTIONS("/adventures/:adventureId/combat", preflight)
	router.OPTIONS("/adventures/:adventureId/:characterId/:op", preflight)

	//Character Endpoints
	//	router.GET("/characters/:campaignId", characterApi.GetCharactersForCampaign)
	//	router.GET("/characters/:adventureId", characterApi.GetCharactersForAdventure)
	//	router.GET("/characters/:characterId", characterApi.GetCharacterById)

	//	router.POST("/characters/:campaignId", characterApi.CreateCharacterForCampaign)

	//	router.PATCH("/characters/:characterId/:attributes", characterApi.UpdateCharacter)

	//	router.DELETE("/characters/:characterId", characterApi.DeleteCharacter)

	router.OPTIONS("/characters", preflight)
	router.OPTIONS("/characters/:adventureId", preflight)

	router.Run("localhost:9090")
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
	c.JSON(http.StatusOK, struct{}{})
}
