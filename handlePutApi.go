package main

import (
	"dict_tagging/funcs"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func handlePut(engine *gin.Engine) {
	engine.POST("/put", func(ctx *gin.Context) {
		dictWriteLock.Lock()
		defer dictWriteLock.Unlock()

		startMicros := time.Now().UnixMicro()
		// 获取上传的文件
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ApiResult{
				Code:   400,
				Msg:    "No file upload",
				Result: "",
				Micros: int(time.Now().UnixMicro() - startMicros),
			})
			return
		}
		// 提取文件名
		baseName := filepath.Base(file.Filename)
		// 分离扩展名
		ext := filepath.Ext(baseName)
		dict := baseName[:len(baseName)-len(ext)]
		ext = strings.ToLower(ext)
		if ext == "" || ext != "json" || dict == "" {
			ctx.JSON(200, ApiResult{
				Code:   100,
				Msg:    "only support *.json",
				Result: "",
				Micros: int(time.Now().UnixMicro() - startMicros),
			})
			return
		}
		dict_path := filepath.Join(config.App.DictDir, fmt.Sprintf("%s.json", dict))
		if funcs.IsPathExist(dict_path) {
			os.Remove(dict_path)
		}
		// 保存文件到指定路径
		if err := ctx.SaveUploadedFile(file, dict_path); err != nil {
			ctx.JSON(http.StatusInternalServerError, ApiResult{
				Code:   500,
				Msg:    fmt.Sprintf("save dict %s failed", dict),
				Result: "",
				Micros: int(time.Now().UnixMicro() - startMicros),
			})
			return
		}
		ctx.JSON(http.StatusOK, ApiResult{
			Code:   1,
			Msg:    "",
			Result: fmt.Sprintf("upload dict %s success", dict),
			Micros: int(time.Now().UnixMicro() - startMicros),
		})
	})
}
