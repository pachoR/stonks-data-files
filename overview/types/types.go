package types

type SymbolOverview struct {
	Description  	string  `json:"description"`
	DisplaySymbol 	string 	`json:"displaySymbol"`
	Symbol 			string 	`json:"symbol"`
	Type 			string 	`json:"type"`
}

type SymbolOverviewMeta struct {
	Count 	int 				`json:"count"`
	Result 	[]SymbolOverview 	`json:"result"`
}