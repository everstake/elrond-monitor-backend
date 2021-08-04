package smodels

type Block struct {
	Hash       string   `json:"hash"`
	Nonce      uint64   `json:"nonce"`
	Shard      uint64   `json:"shard"`
	Epoch      uint64   `json:"epoch"`
	TxCount    uint64   `json:"tx_count"`
	Size       uint64   `json:"size"`
	Proposer   string   `json:"proposer"`
	Miniblocks []string `json:"miniblocks"`
	Timestamp  Time     `json:"timestamp"`
}
