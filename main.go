package main

import (
	"dict_tagging/dict"
	"dict_tagging/statement"
	"fmt"
	"sync"
	"time"
)

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

func search(root *dict.TrieNode, clause string) []dict.QueryEntry {
	startTime := time.Now().UnixMicro()
	searchWords := statement.Split(clause)
	entries := concurrentSearch(root, searchWords)
	endTime := time.Now().UnixMicro()
	fmt.Printf("search cost %d ns\n", (endTime - startTime))
	return entries
}

var (
	// 全局变量
	root *dict.TrieNode
	// 读写锁保护全局变量
	mutex sync.RWMutex
)

func init() {
	root = dict.LoadData()
}

// 获取配置的副本，供协程安全使用
func getRootNode() *dict.TrieNode {
	mutex.RLock()
	defer mutex.RUnlock()
	return root
}

// 异步更新根节点
func asyncUpdateRootNode() {
	// 在独立的协程中加载新数据
	go func() {
		newRoot := dict.LoadData() // 加载新数据，耗时操作

		// 加载完成后，用写锁替换全局的 root
		mutex.Lock()
		defer mutex.Unlock()
		root = newRoot
	}()
}

func main() {
	clause := "我感冒了，医生建议吃对乙酰氨基酚，并且医生嘱咐说：不要吃布洛芬就这样子，另外要调节血糖、注意清咽润喉；也可以弄点奥利司他吃吃"
	entries := search(root, clause)
	for _, entry := range entries {
		fmt.Printf("search %s get entry %+v\n", entry.Dict, entry)
	}
}
