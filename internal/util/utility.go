package util

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func CheckApiErr(e error, ctx *gin.Context) bool {
	if e != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": e.Error()})
		return true
	}
	return false
}

func ApplyCorsHeaders(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
}
