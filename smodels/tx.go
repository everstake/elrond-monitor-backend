package smodels

import "github.com/shopspring/decimal"

type Tx struct {
	Hash          string          `json:"hash"`
	Status        string          `json:"status"`
	From          string          `json:"from"`
	To            string          `json:"to"`
	Value         decimal.Decimal `json:"value"`
	Fee           decimal.Decimal `json:"fee"`
	GasUsed       uint64          `json:"gas_used"`
	MiniblockHash string          `json:"miniblock_hash"`
	ShardFrom     uint64          `json:"shard_from"`
	ShardTo       uint64          `json:"shard_to"`
	Type          string          `json:"type"`
	ScResults     []ScResult      `json:"scResults,omitempty"`
	Timestamp     Time            `json:"timestamp"`
}
