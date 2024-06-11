package api

import (
	"net/http"
	"strconv"

	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
	"github.com/gin-gonic/gin"
)

type CharacterApi struct {
	characterService services.CharacterService
}

func NewCharacterApi(as services.CharacterService) *CharacterApi {
	return &CharacterApi{characterService: as}
}

func (c CharacterApi) CreateCharacterForCampaign(ctx *gin.Context) {
	campaignId, err := strconv.Atoi(ctx.Param("campaignId"))
	campaign := types.NewCampaign(campaignId)
	if util.CheckApiErr(err, ctx) {
		return
	}
	created, err := c.characterService.CreateCharacterForCampaign(campaign)
	if util.CheckApiErr(err, ctx) {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": created})
}

func (c CharacterApi) ManageCharactersForAdventure(ctx *gin.Context) {
	util.ApplyCorsHeaders(ctx)
	adventureId, err := strconv.Atoi(ctx.Param("adventureId"))
	if util.CheckApiErr(err, ctx) {
		return
	}
	characterId, err := strconv.Atoi(ctx.Param("characterId"))
	if util.CheckApiErr(err, ctx) {
		return
	}

	operation, ok := ctx.GetQuery("operation")
	if !ok {
		ctx.JSON(http.StatusBadRequest, util.NamedParameterNotProvided("operation"))
		return
	}
	halfshare := ctx.DefaultQuery("halfshare", "false")

	adventure := types.NewAdventureRecordById(adventureId)
	character := types.NewCharacterById(characterId)
	status, err := c.characterService.ManageCharactersForAdventure(adventure, character, operation, halfshare)
	if util.CheckApiErr(err, ctx) {
		return
	}
	if !status {
		ctx.JSON(http.StatusOK, gin.H{"status": "fail", "data": status})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": status})
}
