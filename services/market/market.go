package market

import (
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/services/market/cmc"
	smodels2 "github.com/everstake/elrond-monitor-backend/smodels"
)

const (
	cmcProvider = "cmc"
)

type (
	Provider interface {
		GetMarketData() (smodels2.MarketData, error)
	}
)

func GetProvider(cfg config.MarketProvider) (p Provider) {
	switch cfg.Title {
	case cmcProvider:
		p = cmc.NewCMC(cfg.APIKey)
	}
	return p
}
