package smodels

type Block struct {
	Hash       string   `json:"hash"`
	Nonce      uint64   `json:"nonce"`
	Shard      uint64   `json:"shard"`
	Epoch      uint64   `json:"epoch"`
	TxCount    uint64   `json:"tx_count"`
	Size       uint64   `json:"size,omitempty"`
	Proposer   string   `json:"propose,omitemptyr"`
	Miniblocks []string `json:"miniblocks,omitempty"`
	Timestamp  Time     `json:"timestamp"`
}
