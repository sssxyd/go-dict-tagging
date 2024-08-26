package main

import (
	"dict_tagging/dict"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func handleReload(engine *gin.Engine) {
	engine.GET("/reload", func(ctx *gin.Context) {
		startMicros := time.Now().UnixMicro()
		dictWriteLock.Lock()
		defer dictWriteLock.Unlock()
		newRoot, newInfos := dict.LoadData()
		root = newRoot
		infos = newInfos
		ctx.JSON(http.StatusOK, ApiResult{
			Code:   1,
			Msg:    "",
			Result: "reload success",
			Micros: int(time.Now().UnixMicro() - startMicros),
		})
	})

	engine.POST("/reload", func(ctx *gin.Context) {
		startMicros := time.Now().UnixMicro()
		dictWriteLock.Lock()
		defer dictWriteLock.Unlock()
		newRoot, newInfos := dict.LoadData()
		root = newRoot
		infos = newInfos
		ctx.JSON(http.StatusOK, ApiResult{
			Code:   1,
			Msg:    "",
			Result: "reload success",
			Micros: int(time.Now().UnixMicro() - startMicros),
		})
	})
}
