package smodels

import "github.com/shopspring/decimal"

type Ranking struct {
	Provider   string                     `json:"provider"`
	Amount     decimal.Decimal            `json:"amount"`
	Delegators map[string]decimal.Decimal `json:"delegators"`
}
