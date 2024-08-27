package main

import (
	"dict_tagging/dict"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// handleReload is a function that handles the "/reload" endpoint for both GET and POST requests.
// It reloads the dictionary data, updates the root and infos variables, and returns a JSON response with the result.
// The function takes in a *gin.Engine parameter, which represents the Gin engine to handle the requests.
// It acquires a write lock on the dictionary using dictWriteLock to ensure thread safety during the reload process.
// After reloading the data, it releases the write lock.
// The function returns the updated root and infos variables and a JSON response with the result and the time taken for the reload operation.
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
