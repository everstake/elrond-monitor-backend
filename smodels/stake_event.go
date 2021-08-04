package smodels

import (
	"github.com/shopspring/decimal"
)

type StakeEvent struct {
	TxHash    string          `json:"tx_hash"`
	Type      string          `json:"type"`
	Validator string          `json:"validator"`
	Delegator string          `json:"delegator"`
	Epoch     uint64          `json:"epoch"`
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt Time            `json:"created_at"`
}
