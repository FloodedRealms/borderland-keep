package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kevin/adventure-archivist/services"
)

type AdventureRecordApi struct {
	adventureRecordService services.AdventureRecordService
}

func NewAdventureRecordApi(as services.AdventureRecordService) *AdventureRecordApi {
	return &AdventureRecordApi{adventureRecordService: as}
}

func (ara AdventureRecordApi) ListAdventureRecordsForCampaign(ctx *gin.Context) {

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
