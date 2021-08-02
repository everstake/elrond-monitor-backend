package smodels

import (
	"github.com/shopspring/decimal"
)

type RangeItem struct {
	Value decimal.Decimal `json:"value"`
	Time  Time            `json:"time"`
}
