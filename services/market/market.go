package market

import (
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/services/market/cmc"
	"github.com/shopspring/decimal"
)

const (
	cmcProvider = "cmc"
)

type (
	Provider interface {
		GetMarketData() (Data, error)
	}
	Data struct {
		Price             decimal.Decimal `json:"price"`
		Cap               decimal.Decimal `json:"cap"`
		CapChange         decimal.Decimal `json:"cap_change"`
		TradingVolume24h  decimal.Decimal `json:"volume_24h"`
		CirculatingSupply decimal.Decimal `json:"circulating_supply"`
		MaxSupply         decimal.Decimal `json:"max_supply"`
		TotalSupply       decimal.Decimal `json:"total_supply"`
	}
)

func GetProvider(cfg config.MarketProvider) (p Provider) {
	switch cfg.Title {
	case cmcProvider:
		p = cmc.NewCMC(cfg.APIKey)
	}
	return p
}
