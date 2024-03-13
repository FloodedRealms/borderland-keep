package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kevin/adventure-archivist/services"
	"github.com/kevin/adventure-archivist/types"
)

type CampaignApi struct {
	campaignService services.CampaignService
}

func NewCampaignApi(cs services.CampaignService) CampaignApi {
	return CampaignApi{campaignService: cs}
}

func (ca *CampaignApi) CreateCampaign(ctx *gin.Context) {
	var cr *types.CreateCampaignRequest

	if err := ctx.ShouldBindJSON(&cr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	nc, err := ca.campaignService.CreateCampaign(cr)
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

func (ca *CampaignApi) ListCampaigns(ctx *gin.Context) {
	applyCorsHeaders(ctx)

	arr, err := ca.campaignService.ListCampaigns()
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
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": campaign})
}

func applyCorsHeaders(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
}
