package smodels

import "github.com/shopspring/decimal"

type ScResult struct {
	Hash    string          `json:"hash"`
	From    string          `json:"from"`
	To      string          `json:"to"`
	Value   decimal.Decimal `json:"value,omitempty"`
	Data    string          `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
}
