package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
	"github.com/gin-gonic/gin"
)

type AdventureRecordApi struct {
	adventureRecordService services.AdventureRecordService
}

func NewAdventureRecordApi(as services.AdventureRecordService) *AdventureRecordApi {
	return &AdventureRecordApi{adventureRecordService: as}
}

func (ara AdventureRecordApi) CreateAdventureRecord(ctx *gin.Context) {
	util.ApplyCorsHeaders(ctx)
	var cr *types.CreateAdventureRecordRequest

	if err := ctx.ShouldBindJSON(&cr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	nc, err := ara.adventureRecordService.CreateAdventureRecordForCampaign(cr)
	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": nc})

}

func (ara AdventureRecordApi) ListAdventureRecordsForCampaign(ctx *gin.Context) {
	util.ApplyCorsHeaders(ctx)

	id := ctx.Param("campaignId")
	arr, err := ara.adventureRecordService.ListAdventureRecordsForCampaign(id)

	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": arr})
}

func (ara AdventureRecordApi) GetAdventure(ctx *gin.Context) {
	util.ApplyCorsHeaders(ctx)
	id := ctx.Param("adventureId")
	arr, err := ara.adventureRecordService.GetAdventureRecordById(id)

	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": arr})

}

func (ara AdventureRecordApi) AddLootToAdventure(ctx *gin.Context) {
	log.Print("hit")
	util.ApplyCorsHeaders(ctx)
	adventureId, err := strconv.Atoi(ctx.Param("adventureId"))
	util.CheckErr(err)
	adventure := types.NewAdventureRecordById(adventureId)
	lootType := ctx.Query("type")
	var lootObject *types.XPSource
	if err = ctx.ShouldBindJSON(&lootObject); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log.Print(lootObject)
	if lootType == "" {
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": "fail", "message": util.RequiredParameterNotProvided()})
		return
	}
	switch lootType {
	case string(types.GemLoot):
		gem := types.NewGem(lootObject.Name, lootObject.Description, lootObject.XPValue, lootObject.Number)
		log.Print(ctx.Request.Body)
		ara.adventureRecordService.AddGemLootToAdventure(adventure, gem)
	default:
		ctx.JSON(http.StatusNotImplemented, util.NotYetImplmented())

	}
}
