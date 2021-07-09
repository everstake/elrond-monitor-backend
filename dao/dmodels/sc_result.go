package dmodels

import (
	"github.com/shopspring/decimal"
)

const SCResultsTable = "sc_results"

type SCResult struct {
	Hash    string          `db:"scr_hash"`
	TxHash  string          `db:"trn_hash"`
	From    string          `db:"scr_from"`
	To      string          `db:"scr_to"`
	Value   decimal.Decimal `db:"scr_value"`
	Data    string          `db:"scr_data"`
	Message string          `json:"scr_message"`
}
