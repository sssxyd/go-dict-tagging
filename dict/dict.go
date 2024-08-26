package dict

import (
	"bufio"
	"dict_tagging/funcs"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type WordEntry struct {
	Dict  string
	Word  string                 `json:"word"`
	Index []string               `json:"index"`
	Data  map[string]interface{} `json:"data"`
}

type TrieNode struct {
	Children map[rune]*TrieNode
	Entrys   []WordEntry
}

func listDicts(dirPath string) []string {
	// 读取目录
	dirs, err := os.ReadDir(dirPath)
	if err != nil {
		return make([]string, 0)
	}

	// 过滤出 .dict 后缀的文件
	var dicts []string
	for _, dir := range dirs {
		if !dir.IsDir() {
			fileName := dir.Name()
			if strings.HasSuffix(strings.ToLower(fileName), ".json") {
				dicts = append(dicts, fileName[:len(fileName)-5])
			}
		}
	}
	return dicts
}

func readDictFile(dict string, filePath string) []WordEntry {
	var entries []WordEntry

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error when open dict file: %s, %v\n", filePath, err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		line = strings.TrimRight(line, ",")
		if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
			continue
		}

		var currentEntry WordEntry
		err := json.Unmarshal([]byte(line), &currentEntry)
		if err != nil {
			log.Fatalf("error parsing dict file %s: %v\n", filePath, err)
		}
		currentEntry.Dict = dict
		entries = append(entries, currentEntry)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error when scanning file %s: %v\n", filePath, err)
	}

	return entries
}

func insertIntoTrie(root *TrieNode, entry WordEntry) {
	for _, index := range entry.Index {
		currentNode := root
		for _, ch := range index {
			if currentNode.Children == nil {
				currentNode.Children = make(map[rune]*TrieNode)
			}
			if _, exists := currentNode.Children[ch]; !exists {
				currentNode.Children[ch] = &TrieNode{
					Children: make(map[rune]*TrieNode),
				}
			}
			currentNode = currentNode.Children[ch]
		}
		currentNode.Entrys = append(currentNode.Entrys, entry)
	}
}

func LoadData() *TrieNode {
	startTime := time.Now().UnixMilli()
	dict_path := filepath.Join(funcs.GetExecutionPath(), "data")
	dict_names := listDicts(dict_path)
	fmt.Printf("load dict %v from path %s\n", dict_names, dict_path)
	if len(dict_names) == 0 {
		return &TrieNode{}
	}
	root := TrieNode{
		Children: make(map[rune]*TrieNode),
	}
	total := 0
	for _, dict := range dict_names {
		fmt.Printf("read dict %s\n", dict)
		wordEnties := readDictFile(dict, filepath.Join(dict_path, dict+".json"))
		if wordEnties == nil {
			continue
		}
		total += len(wordEnties)
		for _, entry := range wordEnties {
			insertIntoTrie(&root, entry)
		}
	}
	endTime := time.Now().UnixMilli()
	fmt.Printf("load [%s] dicts from data, total %d word entries, cost %d ms\n", strings.Join(dict_names, ","), total, (endTime - startTime))
	return &root
}
