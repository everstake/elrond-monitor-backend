package swagger

type PaginatableResponse struct {
	Count float64 `json:"count,omitempty"`

	Items interface{} `json:"items,omitempty"`
	//Blocks []Block `json:"blocks,omitempty"`
	//Txs []Tx `json:"txs,omitempty"`
	//Validators []Validator `json:"validators,omitempty"`
	//Accounts []Account `json:"accounts,omitempty"`
}
