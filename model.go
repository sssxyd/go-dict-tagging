package main

type PosRange struct {
	start uint16
	end   uint16
}

type StatementEntry struct {
	Pos  []PosRange             `json:"pos"`
	Dict string                 `json:"dict"`
	Word string                 `json:"word"`
	Data map[string]interface{} `json:"data"`
}

type StatementTagResult struct {
	Statement string                      `json:"statement"`
	DictTags  map[string][]StatementEntry `json:"dictTags"`
	CostMicro uint64                      `json:"costMicro"`
}
