package smodels

import "github.com/shopspring/decimal"

type MarketData struct {
	Price             decimal.Decimal `json:"price"`
	Cap               decimal.Decimal `json:"cap"`
	CapChange         decimal.Decimal `json:"cap_change"`
	TradingVolume24h  decimal.Decimal `json:"volume_24h"`
	CirculatingSupply decimal.Decimal `json:"circulating_supply"`
	MaxSupply         decimal.Decimal `json:"max_supply"`
	TotalSupply       decimal.Decimal `json:"total_supply"`
}
