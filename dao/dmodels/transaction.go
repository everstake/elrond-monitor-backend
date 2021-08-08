package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	TransactionsTable = "transactions"

	TxStatusPending = "pending"
	TxStatusSuccess = "success"
	TxStatusFail    = "fail"
	TxStatusInvalid = "invalid"
)

type Transaction struct {
	Hash          string          `db:"trn_hash"`
	Status        string          `db:"trn_status"`
	MiniBlockHash string          `db:"mlk_mini_block_hash"`
	Value         decimal.Decimal `db:"trn_value"`
	Sender        string          `db:"trn_sender"`
	SenderShard   uint64          `db:"trn_sender_shard"`
	Receiver      string          `db:"trn_receiver"`
	ReceiverShard uint64          `db:"trn_receiver_shard"`
	GasPrice      uint64          `db:"trn_gas_price"`
	Nonce         uint64          `db:"trn_nonce"`
	Data          string          `db:"trn_data"`
	CreatedAt     time.Time       `db:"trn_created_at"`
}
