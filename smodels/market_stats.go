package smodels

type MarketStats struct {
	MarketCap float64 `json:"market_cap,omitempty"`

	CirculatingSupply float64 `json:"circulating_supply,omitempty"`

	MaxPrice float64 `json:"max_price,omitempty"`

	Price float64 `json:"price,omitempty"`
}
