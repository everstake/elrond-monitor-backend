package smodels

import "github.com/shopspring/decimal"

type Stats struct {
	Price             decimal.Decimal `json:"price"`
	PriceChange       decimal.Decimal `json:"price_change"`
	TradingVolume     decimal.Decimal `json:"trading_volume"`
	Cap               decimal.Decimal `json:"cap"`
	CapChange         decimal.Decimal `json:"cap_change"`
	CirculatingSupply decimal.Decimal `json:"circulating_supply"`
	TotalSupply       decimal.Decimal `json:"total_supply"`
	Height            uint64          `json:"height"`
	TotalTxs          uint64          `json:"total_txs"`
	TotalAccounts     uint64          `json:"total_accounts"`
}

type ValidatorStats struct {
	ActiveStake   decimal.Decimal `json:"active_stake"`
	Validators    uint64          `json:"validators"`
	ObserverNodes uint64          `json:"observer_nodes"`
	StakingAPR    decimal.Decimal `json:"staking_apr"`
	Queue         uint64          `json:"queue"`
}
