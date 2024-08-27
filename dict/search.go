package dict

type QueryEntry struct {
	Start uint16
	End   uint16
	Entry WordEntry
}

func Search(root *TrieNode, keyword string, start uint16) []QueryEntry {
	currentNode := root
	var queryResults []QueryEntry
	var index uint16
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
				queryResults = append(queryResults, QueryEntry{
					Start: start,
					End:   start + index,
					Entry: entry,
				})
			}
		}
		index++
	}
	return queryResults
}
