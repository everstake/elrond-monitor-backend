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
	Data          string          `json:"data"`
	Timestamp     Time            `json:"timestamp"`
}

type Operation struct {
	Nonce          uint64            `json:"nonce"`
	Sender         string            `json:"sender"`
	Receiver       string            `json:"receiver"`
	OriginalTxHash string            `json:"original_tx_hash"`
	Timestamp      uint64            `json:"timestamp"`
	Status         string            `json:"status"`
	SenderShard    uint64            `json:"sender_shard"`
	ReceiverShard  uint64            `json:"receiver_shard"`
	Operation      string            `json:"operation"`
	Tokens         []string          `json:"tokens"`
	ESDTValues     []decimal.Decimal `json:"esdt_values"`
}
