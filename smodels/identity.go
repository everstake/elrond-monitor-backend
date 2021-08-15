package smodels

import "github.com/shopspring/decimal"

type Identity struct {
	Avatar       string          `json:"avatar"`
	Description  string          `json:"description"`
	Identity     string          `json:"identity"`
	Locked       decimal.Decimal `json:"locked"`
	Name         string          `json:"name"`
	Rank         uint64          `json:"rank"`
	Score        uint64          `json:"score"`
	Stake        decimal.Decimal `json:"stake"`
	StakePercent float64         `json:"stake_percent"`
	ToUp         decimal.Decimal `json:"to_up"`
	Validators   uint64          `json:"validators"`
}
