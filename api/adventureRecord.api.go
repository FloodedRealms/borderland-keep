package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/types"
)

type AdventureApi struct {
	adventureRecordService services.AdventureService
	characterRecordService services.CharacterService
}

func NewAdventureRecordApi(as services.AdventureService, cs services.CharacterService) *AdventureApi {
	return &AdventureApi{adventureRecordService: as, characterRecordService: cs}
}

func (ara AdventureApi) CreateAdventureRecord(w http.ResponseWriter, r *http.Request) {
	var cr types.AdventureRecord
	err := decodeJSONBody(w, r, &cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nc, err := ara.adventureRecordService.CreateAdventureRecordForCampaign(&cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendGoodResponseWithObject(w, nc)
}

// TODO: Fix the N+1 selection here
func (ara AdventureApi) ListAdventureRecordsForCampaign(w http.ResponseWriter, r *http.Request) {
	applyCorsHeaders(w)
	id := r.PathValue("campaignId")
	arr, err := ara.adventureRecordService.ListAdventureRecordsForCampaign(id)

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendGoodResponseWithObject(w, arr)
}

func (ara AdventureApi) GetAdventure(w http.ResponseWriter, r *http.Request) {
	i := r.PathValue("adventureId")
	id, _ := strconv.Atoi(i)
	arr, err := ara.adventureRecordService.GetAdventureRecordById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendGoodResponseWithObject(w, arr)
}

func (ara AdventureApi) UpdateAdventure(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("adventureId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var ur types.AdventureRecord
	err = decodeJSONBody(w, r, &ur)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ur.Id = id
	uAdventure, err := ara.adventureRecordService.UpdateAdventureRecord(&ur)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendGoodResponseWithObject(w, uAdventure)

}

/* func (ara AdventureRecordApi) AddLootToAdventure(w http.ResponseWriter, r *http.Request) {
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

func (ara AdventureApi) makeBoolResponse(status bool, err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if status != true {
		http.Error(w, errors.New("operation failed").Error(), http.StatusExpectationFailed)
		return
	}
	sendSuccessResponse(w, "Operation Succeeded")

}
