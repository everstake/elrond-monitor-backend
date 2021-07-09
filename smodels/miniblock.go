package smodels

type Miniblock struct {
	Hash          string `json:"hash"`
	ShardFrom     uint64 `json:"shard_from"`
	ShardTo       uint64 `json:"shard_to"`
	BlockSender   string `json:"block_sender"`
	BlockReceiver string `json:"block_receiver"`
	Type          string `json:"type,omitempty"`
	Txs           []Tx   `json:"txs,omitempty"`
	Timestamp     Time   `json:"timestamp"`
}
