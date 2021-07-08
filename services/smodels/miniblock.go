package swagger

type Miniblock struct {
	Hash string `json:"hash,omitempty"`

	ShardFrom string `json:"shard_from,omitempty"`

	ShardTo string `json:"shard_to,omitempty"`

	BlockSender string `json:"block_sender,omitempty"`

	BlockReceiver string `json:"block_receiver,omitempty"`

	Type string `json:"type,omitempty"`

	Txs []Tx `json:"txs,omitempty"`

	Timestamp float64 `json:"timestamp,omitempty"`
}
