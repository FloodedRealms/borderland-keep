package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kevin/adventure-archivist/services"
	"github.com/kevin/adventure-archivist/types"
	"github.com/kevin/adventure-archivist/util"
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
