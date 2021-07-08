package dmodels

import (
	"github.com/shopspring/decimal"
)

const (
	DelegationsTable = "delegations"
	StakesTable      = "stakes"
)

type (
	Delegation struct {
		ID        string          `db:"dlg_id"`
		TxHash    string          `db:"dlg_tx_hash"`
		Delegator string          `db:"dlg_delegator"`
		Validator string          `db:"dlg_validator"`
		Amount    decimal.Decimal `db:"dlg_amount"`
		CreatedAt Time            `db:"dlg_created_at"`
	}

	Stake struct {
		ID        string          `db:"stk_id"`
		TxHash    string          `db:"stk_tx_hash"`
		Delegator string          `db:"stk_delegator"`
		Validator string          `db:"stk_validator"`
		Amount    decimal.Decimal `db:"stk_amount"`
		CreatedAt Time            `db:"stk_created_at"`
	}
)
