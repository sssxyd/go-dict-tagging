package dict

type QueryEntry struct {
	Start uint16                 `json:"start"`
	End   uint16                 `json:"end"`
	Dict  string                 `json:"dict"`
	Index []string               `json:"index"`
	Word  string                 `json:"word"`
	Data  map[string]interface{} `json:"data"`
}

func Search(root *TrieNode, keyword string, start uint16) []QueryEntry {
	currentNode := root
	var queryResults []QueryEntry
	var index uint16
	entrySet := make(map[string]bool)
	for _, ch := range keyword {
		if currentNode.Children == nil {
			break
		}
		nextNode, exists := currentNode.Children[ch]
		if !exists {
			break
		}
		currentNode = nextNode
		if len(currentNode.Entrys) > 0 {
			for _, entry := range currentNode.Entrys {
				entry_key := entry.Dict + "_" + entry.Word
				if _, exists := entrySet[entry_key]; !exists {
					entrySet[entry_key] = true
					queryResults = append(queryResults, QueryEntry{
						Start: start,
						End:   start + index,
						Dict:  entry.Dict,
						Index: entry.Index,
						Word:  entry.Word,
						Data:  entry.Data,
					})
				}
			}
		}
		index++
	}
	return queryResults
}
