package swagger

type Tx struct {
	Hash string `json:"hash,omitempty"`

	Status string `json:"status,omitempty"`

	From string `json:"from,omitempty"`

	To string `json:"to,omitempty"`

	Value float64 `json:"value,omitempty"`

	Fee float64 `json:"fee,omitempty"`

	GasUsed float64 `json:"gas_used,omitempty"`

	MiniblockHash string `json:"miniblock_hash,omitempty"`

	ShardFrom string `json:"shard_from,omitempty"`

	ShardTo string `json:"shard_to,omitempty"`

	Type string `json:"type,omitempty"`

	ScResults []ScResult `json:"scResults,omitempty"`

	Timestamp float64 `json:"timestamp,omitempty"`
}
