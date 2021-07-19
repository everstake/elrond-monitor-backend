package smodels

import "github.com/shopspring/decimal"

type Validator struct {
	Identity     string                 `json:"identity"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Avatar       string                 `json:"avatar"`
	Score        uint64                 `json:"score"`
	Validators   uint64                 `json:"validators"`
	Stake        string                 `json:"stake"`
	TopUp        string                 `json:"topUp"`
	Locked       string                 `json:"locked"`
	Distribution map[string]interface{} `json:"distribution"`
	Providers    []string               `json:"providers"`
	StakePercent decimal.Decimal        `json:"stakePercent"`
	Rank         uint64                 `json:"rank"`
}
