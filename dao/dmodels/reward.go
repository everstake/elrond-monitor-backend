package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const RewardsTable = "rewards"

type Reward struct {
	TxHash          string          `db:"rwd_tx_hash"`
	HypeblockID     uint64          `db:"rwd_hyperblock_id"`
	ReceiverAddress string          `db:"rwd_receiver_address"`
	Amount          decimal.Decimal `db:"rwd_amount"`
	CreatedAt       time.Time       `db:"rwd_created_at"`
}
