package swagger

type Block struct {
	Hash string `json:"hash,omitempty"`

	Nonce float64 `json:"nonce,omitempty"`

	Shard string `json:"shard,omitempty"`

	Epoch float64 `json:"epoch,omitempty"`

	TxCount string `json:"tx_count,omitempty"`

	Size float64 `json:"size,omitempty"`

	Proposer string `json:"proposer,omitempty"`

	Miniblocks []string `json:"miniblocks,omitempty"`

	Timestamp float64 `json:"timestamp,omitempty"`
}
