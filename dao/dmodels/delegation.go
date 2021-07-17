package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	DelegationsTable = "delegations"
	StakesTable      = "stakes"
)

type (
	Delegation struct {
		TxHash    string          `db:"dlg_tx_hash"`
		Delegator string          `db:"dlg_delegator"`
		Validator string          `db:"dlg_validator"`
		Amount    decimal.Decimal `db:"dlg_amount"`
		CreatedAt time.Time       `db:"dlg_created_at"`
	}

	Stake struct {
		TxHash    string          `db:"stk_tx_hash"`
		Validator string          `db:"stk_validator"`
		Amount    decimal.Decimal `db:"stk_amount"`
		CreatedAt time.Time       `db:"stk_created_at"`
	}
)
