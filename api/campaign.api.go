package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/gin-gonic/gin"
)

type CampaignApi struct {
	campaignService  services.CampaignService
	characterService services.CharacterService
}

func NewCampaignApi(cs services.CampaignService, chars services.CharacterService) CampaignApi {
	return CampaignApi{
		campaignService:  cs,
		characterService: chars}
}

func (ca *CampaignApi) CreateCampaign(ctx *gin.Context) {
	var cr *types.CreateCampaignRecordRequest
	clientId := ctx.Request.Header.Get("X-Archivist-Client-Id")
	if err := ctx.ShouldBindJSON(&cr); err != nil {
		body, _ := io.ReadAll(ctx.Request.Body)
		println(string(body))
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	nc, err := ca.campaignService.CreateCampaign(cr, clientId)
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

func (ca *CampaignApi) UpdateCampaign(ctx *gin.Context) {
	var cr *types.CampaignRecord

	if err := ctx.ShouldBindJSON(&cr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	nc, err := ca.campaignService.UpdateCampaign(cr)
	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": nc})

}

func (ca *CampaignApi) ListCampaigns(ctx *gin.Context) {
	applyCorsHeaders(ctx)
	clientId := ctx.Request.Header.Get("X-Archivist-Client-Id")

	var arr []*types.CampaignRecord
	var err error
	if clientId != "" {

	} else {
		arr, err = ca.campaignService.ListCampaigns()
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": arr})
}

func (ca *CampaignApi) GetCampaign(ctx *gin.Context) {
	applyCorsHeaders(ctx)
	id := ctx.Param("campaignId")

	campaign, err := ca.campaignService.GetCampaign(id)

	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	campaign.Characters, err = ca.characterService.GetCharactersForCampaign(campaign)
	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": campaign})
}

func (ca *CampaignApi) DeleteCampaign(ctx *gin.Context) {
	applyCorsHeaders(ctx)
	id := ctx.Param("campaignId")

	campaign, err := ca.campaignService.DeleteCampaign(id)

	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": campaign})
}

func applyCorsHeaders(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
}
