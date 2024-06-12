package api

import (
	"net/http"
	"strings"

	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
	"github.com/gin-gonic/gin"
)

type AdventureApi struct {
	adventureRecordService services.AdventureService
}

func NewAdventureRecordApi(as services.AdventureService) *AdventureApi {
	return &AdventureApi{adventureRecordService: as}
}

func (ara AdventureApi) CreateAdventureRecord(ctx *gin.Context) {
	util.ApplyCorsHeaders(ctx)
	var cr *types.CreateAdventureRequest

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

func (ara AdventureApi) ListAdventureRecordsForCampaign(ctx *gin.Context) {
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

func (ara AdventureApi) GetAdventure(ctx *gin.Context) {
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

func (ara AdventureApi) UpdateAdventure(ctx *gin.Context) {
	util.ApplyCorsHeaders(ctx)
	var ur *types.UpdateAdventureRequest

	if err := ctx.ShouldBindJSON(&ur); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	uAdventure, err := ara.adventureRecordService.UpdateAdventureRecord(ur)
	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not yet implemented") {
			ctx.JSON(http.StatusNotImplemented, gin.H{"status": "fail", "message": ur})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": uAdventure})

}

/* func (ara AdventureRecordApi) AddLootToAdventure(ctx *gin.Context) {
	log.Print("hit")
	util.ApplyCorsHeaders(ctx)
	adventureId, err := strconv.Atoi(ctx.Param("adventureId"))
	util.CheckErr(err)
	adventure := types.NewAdventureRecordById(adventureId)
	lootType := ctx.Query("type")

	if lootType == "" {
		ctx.JSON(http.StatusPreconditionFailed, gin.H{"status": "fail", "message": util.RequiredParameterNotProvided()})
		return
	}
	switch lootType {
	case string(types.CoinLoot):
		var coinObject *types.CoinUpdateRequest
		if err = ctx.ShouldBindJSON(&coinObject); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		coins := types.NewCoins(coinObject.Copper, coinObject.Silver, coinObject.Electrum, coinObject.Gold, coinObject.Platinum)
		status, err := ara.adventureRecordService.AddCoinsAdventure(adventure, coins)
		ara.makeBoolResponse(status, err, ctx)
	case string(types.GemLoot):
		var lootObject *types.XPSource
		if err = ctx.ShouldBindJSON(&lootObject); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		gem := types.NewGem(lootObject.Name, lootObject.Description, lootObject.XPValue, lootObject.Number, -1)
		log.Print(ctx.Request.Body)
		status, err := ara.adventureRecordService.AddGemLootToAdventure(adventure, gem)
		ara.makeBoolResponse(status, err, ctx)

	case string(types.JewelleryLoot):
		var lootObject *types.XPSource
		if err = ctx.ShouldBindJSON(&lootObject); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		jewellery := types.NewJewellery(lootObject.Name, lootObject.Description, lootObject.XPValue, lootObject.Number, -1)
		log.Print(ctx.Request.Body)
		status, err := ara.adventureRecordService.AddJewelleryLootToAdventure(adventure, jewellery)
		ara.makeBoolResponse(status, err, ctx)

	case string(types.MagicItemLoot):
		var lootObject *types.MagicItemRequest
		if err = ctx.ShouldBindJSON(&lootObject); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		magicItem := types.NewMagicItem(lootObject.Name, lootObject.Description, float64(lootObject.ApparentValue), lootObject.ActualValue, -1)
		status, err := ara.adventureRecordService.AddMagicItemToAdventure(adventure, magicItem)
		ara.makeBoolResponse(status, err, ctx)

	case string(types.Combat):
		var lootObject *types.MonsterGroupRequest
		if err = ctx.ShouldBindJSON(&lootObject); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		combat := types.NewMonsterGroup(lootObject.MonsterName, lootObject.NumberDefeated, -1, float64(lootObject.XPPerMonster))
		status, err := ara.adventureRecordService.AddCombatToAdventure(adventure, combat)
		ara.makeBoolResponse(status, err, ctx)

	default:
		ctx.JSON(http.StatusNotImplemented, util.NotYetImplmented())

	}
} */

func (ara AdventureApi) makeBoolResponse(status bool, err error, ctx *gin.Context) {
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	if status != true {
		ctx.JSON(http.StatusExpectationFailed, gin.H{"status": "fail", "message": "The requested action failed to complete, but no error was report"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "The requested action succeeded"})

}
