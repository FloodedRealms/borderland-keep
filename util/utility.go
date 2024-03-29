package util

import (
	"log"

	"github.com/gin-gonic/gin"
)

func CheckErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func ApplyCorsHeaders(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
}
