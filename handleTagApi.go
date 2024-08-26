package main

import (
	"dict_tagging/dict"
	"dict_tagging/statement"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Position struct {
	Start uint16 `json:"start"`
	End   uint16 `json:"end"`
}

type KeywordTag struct {
	Keyword   string                      `json:"keyword"`
	Positions []Position                  `json:"positions"`
	DictWords map[string][]dict.WordEntry `json:"dictwords"`
}

type RequestParams struct {
	Statement string `json:"statement" form:"statement"`
}

func (tag *KeywordTag) addQueryWord(query dict.QueryEntry) {
	// 检查是否已经有相同的开始位置
	exists := false
	for _, p := range tag.Positions {
		if p.Start == query.Start {
			exists = true
			break
		}
	}
	if !exists {
		tag.Positions = append(tag.Positions, Position{query.Start, query.End})
	}

	// 添加字典词条
	if dictEntries, exists := tag.DictWords[query.Entry.Dict]; exists {
		alreadyExists := false
		for _, de := range dictEntries {
			if de.Word == query.Entry.Word {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			tag.DictWords[query.Entry.Dict] = append(tag.DictWords[query.Entry.Dict], query.Entry)
		}
	} else {
		// 如果字典键不存在，初始化并添加新的词条
		tag.DictWords[query.Entry.Dict] = append(tag.DictWords[query.Entry.Dict], query.Entry)
	}
}

func fillKeywordTags(statement string, queryEntries []dict.QueryEntry) map[string]KeywordTag {
	keywordTags := make(map[string]KeywordTag)
	runes := []rune(statement)
	for _, qe := range queryEntries {
		keyword := string(runes[qe.Start : qe.End+1])
		if _, exists := keywordTags[keyword]; !exists {
			keywordTags[keyword] = KeywordTag{
				Keyword:   keyword,
				DictWords: make(map[string][]dict.WordEntry),
			}
		}
		// 使用引用来获取和更新 map 中的值
		tag := keywordTags[keyword]
		tag.addQueryWord(qe)
		// 重新赋值回 map 中
		keywordTags[keyword] = tag
	}
	return keywordTags
}

// 并发搜索多个关键词并合并结果
func concurrentSearch(root *dict.TrieNode, searchWords []statement.SearchWord) []dict.QueryEntry {
	var wg sync.WaitGroup
	resultChan := make(chan []dict.QueryEntry, len(searchWords))
	// 对每个关键词启动一个 goroutine
	for _, kw := range searchWords {
		wg.Add(1)
		go func(kw statement.SearchWord) {
			defer wg.Done()
			// 调用 Search 函数
			for _, idx := range kw.Start {
				results := dict.Search(root, kw.Word, idx)
				resultChan <- results
			}
		}(kw)
	}

	// 等待所有的 goroutine 完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 合并所有结果
	var mergedResults []dict.QueryEntry
	for res := range resultChan {
		mergedResults = append(mergedResults, res...)
	}

	return mergedResults
}

func search(clause string) map[string]KeywordTag {
	searchWords := statement.Split(clause)
	entries := concurrentSearch(root, searchWords)
	tags := fillKeywordTags(clause, entries)
	return tags
}

func handleSearch(ctx *gin.Context, startMicros int64, clause string) {
	if clause == "" {
		ctx.JSON(200, ApiResult{
			Code:   100,
			Msg:    "statement is empty",
			Result: "",
			Micros: int(time.Now().UnixMicro() - startMicros),
		})
	} else {
		result := search(clause)
		ctx.JSON(200, ApiResult{
			Code:   1,
			Msg:    "",
			Result: result,
			Micros: int(time.Now().UnixMicro() - startMicros),
		})
	}
}

func handleTag(engine *gin.Engine) {
	engine.GET("/tag", func(ctx *gin.Context) {
		startMicros := time.Now().UnixMicro()
		statement := ctx.DefaultQuery("statement", "")
		handleSearch(ctx, startMicros, statement)
	})

	engine.POST("/tag", func(ctx *gin.Context) {
		startMicros := time.Now().UnixMicro()
		var params RequestParams

		// 获取Content-Type
		contentType := ctx.GetHeader("Content-Type")

		// 根据Content-Type判定
		if contentType == "application/json" {
			// 处理JSON数据
			if err := ctx.ShouldBindJSON(&params); err != nil {
				ctx.JSON(http.StatusBadRequest, ApiResult{
					Code:   400,
					Msg:    err.Error(),
					Result: "",
				})
				return
			}
		} else if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
			// 处理表单数据
			if err := ctx.ShouldBind(&params); err != nil {
				ctx.JSON(http.StatusBadRequest, ApiResult{
					Code:   400,
					Msg:    err.Error(),
					Result: "",
				})
				return
			}
		} else {
			ctx.JSON(http.StatusUnsupportedMediaType, ApiResult{
				Code:   415,
				Msg:    fmt.Sprintf("contentType: %s not supported", contentType),
				Result: "",
			})
			return
		}
		handleSearch(ctx, startMicros, params.Statement)
	})

}
