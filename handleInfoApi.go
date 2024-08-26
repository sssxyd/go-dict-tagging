package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func handleInfo(engine *gin.Engine) {
	engine.GET("/info", func(ctx *gin.Context) {
		startMicros := time.Now().UnixMicro()
		ctx.JSON(http.StatusOK, ApiResult{
			Code:   1,
			Msg:    "",
			Result: infos,
			Micros: int(time.Now().UnixMicro() - startMicros),
		})
	})

	engine.POST("/info", func(ctx *gin.Context) {
		startMicros := time.Now().UnixMicro()
		ctx.JSON(http.StatusOK, ApiResult{
			Code:   1,
			Msg:    "",
			Result: infos,
			Micros: int(time.Now().UnixMicro() - startMicros),
		})
	})
}
