package smodels

type Block struct {
	Hash                  string   `json:"hash"`
	Nonce                 uint64   `json:"nonce"`
	Shard                 uint64   `json:"shard"`
	Epoch                 uint64   `json:"epoch"`
	TxCount               uint64   `json:"tx_count"`
	Size                  int64    `json:"size"`
	Proposer              string   `json:"proposer"`
	Miniblocks            []string `json:"miniblocks"`
	NotarizedBlocksHashes []string `json:"notarizedBlocksHashes"`
	Validators            []string `json:"validators"`
	PubKeyBitmap          string   `json:"pubKeyBitmap"`
	StateRootHash         string   `json:"stateRootHash"`
	PrevHash              string   `json:"prevHash"`
	Timestamp             Time     `json:"timestamp"`
}
