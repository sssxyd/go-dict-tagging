package statement

type index_word struct {
	start uint16
	word  string
}

type SearchWord struct {
	Start []uint16
	Word  string
}

func Split(statement string) []SearchWord {
	var words []index_word
	// 不再做切词，完全匹配，包括空格和标点符号
	// var lastWord string
	// var startIdx uint16
	// runes := []rune(statement)
	// for i, c := range runes {
	// 	if funcs.RuneIsStopChar(c) {
	// 		if lastWord != "" {
	// 			words = append(words, index_word{start: startIdx, word: lastWord})
	// 		}
	// 		// 更新startIdx为下一个词的起始位置
	// 		startIdx = uint16(i + 1)
	// 		lastWord = ""
	// 		continue
	// 	}
	// 	lastWord += string(c)
	// }
	// // 如果最后一个词存在，添加到words中
	// if lastWord != "" {
	// 	words = append(words, index_word{start: startIdx, word: lastWord})
	// }

	words = append(words, index_word{start: 0, word: statement})

	var allWords []index_word
	for _, w := range words {
		runes := []rune(w.word)
		length := len(runes)
		for i := 0; i < length-1; i++ {
			allWords = append(allWords, index_word{
				start: w.start + uint16(i),
				word:  string(runes[i:]),
			})
		}
	}

	var searchWordList []SearchWord
	searchWordMap := make(map[string]int)
	for _, w := range allWords {
		idx, ok := searchWordMap[w.word]
		if !ok {
			sw := SearchWord{
				Start: []uint16{w.start},
				Word:  w.word,
			}
			searchWordList = append(searchWordList, sw)
			searchWordMap[w.word] = len(searchWordList) - 1
		} else {
			searchWordList[idx].Start = append(searchWordList[idx].Start, w.start)
		}
	}

	return searchWordList
}
