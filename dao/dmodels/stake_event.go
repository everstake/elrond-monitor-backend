package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	StakeEventsTable = "stake_events"

	ClaimRewardsEventType      = "claimRewards"
	DelegateStakeEventType     = "delegate"
	UnDelegateStakeEventType   = "unDelegate"
	ReDelegateRewardsEventType = "reDelegateRewards"
	WithdrawEventType          = "withdraw"
	StakeStakeEventType        = "stake"
	UnStakeEventType           = "unStake"
	ReStakeRewardsEventType    = "reStakeRewards"
	UnBondEventType            = "unBond"
)

type StakeEvent struct {
	TxHash    string          `db:"ste_tx_hash"`
	Type      string          `db:"ste_type"`
	Validator string          `db:"ste_validator"`
	Delegator string          `db:"ste_delegator"`
	Epoch     uint64          `db:"ste_epoch"`
	Amount    decimal.Decimal `db:"ste_amount"`
	CreatedAt time.Time       `db:"ste_created_at"`
}
