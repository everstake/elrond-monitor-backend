package dmodels

import (
	"github.com/shopspring/decimal"
)

const RewardsTable = "rewards"

type Reward struct {
	ID              string          `db:"rwd_id"`
	HypeblockID     uint64          `db:"rwd_hyperblock_id"`
	TxHash          string          `db:"rwd_tx_hash"`
	ReceiverAddress string          `db:"rwd_receiver_address"`
	Amount          decimal.Decimal `db:"rwd_amount"`
	CreatedAt       Time            `db:"rwd_created_at"`
}
