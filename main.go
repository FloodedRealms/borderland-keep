package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kevin/adventure-archivist/api"
	"github.com/kevin/adventure-archivist/repository"
	"github.com/kevin/adventure-archivist/services"
	"github.com/kevin/adventure-archivist/util"
)

func main() {
	repo, err := repository.NewMemoryRepo()
	util.CheckErr(err)
	campaignService := services.NewCampaignService(repo, context.TODO())
	campaignApi := api.NewCampaignApi(campaignService)
	adventurRecordService := services.NewAdventureRecordService(repo, context.TODO())
	adventurRecordApi := api.NewAdventureRecordApi(adventurRecordService)
	router := gin.Default()
	router.POST("/campaigns", campaignApi.CreateCampaign)
	router.GET("/campaigns", campaignApi.ListCampaigns)
	router.GET("/campaigns/:campaignId", campaignApi.GetCampaign)
	router.GET("/campaigns/:campaignId/adventures", adventurRecordApi.ListAdventureRecordsForCampaign)
	router.POST("/adventures", adventurRecordApi.CreateAdventureRecord)

	router.OPTIONS("/campaigns", preflight)
	router.OPTIONS("/campaigns/:campaignId", preflight)
	router.OPTIONS("/adventures", preflight)
	router.Run("localhost:9090")
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
	c.JSON(http.StatusOK, struct{}{})
}
