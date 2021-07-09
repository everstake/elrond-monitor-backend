package smodels

type Block struct {
	Hash       string   `json:"hash"`
	Nonce      float64  `json:"nonce"`
	Shard      string   `json:"shard"`
	Epoch      float64  `json:"epoch"`
	TxCount    string   `json:"tx_count"`
	Size       float64  `json:"size"`
	Proposer   string   `json:"proposer"`
	Miniblocks []string `json:"miniblocks,omitempty"`
	Timestamp  float64  `json:"timestamp"`
}
