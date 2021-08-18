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
	GasPrice      uint64          `json:"gas_price"`
	MiniblockHash string          `json:"miniblock_hash"`
	ShardFrom     uint64          `json:"shard_from"`
	ShardTo       uint64          `json:"shard_to"`
	ScResults     []ScResult      `json:"scResults"`
	Signature     string          `json:"signature"`
	Timestamp     Time            `json:"timestamp"`
}
