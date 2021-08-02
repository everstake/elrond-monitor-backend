package market

import (
	"errors"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/services/market/cmc"
	"github.com/everstake/elrond-monitor-backend/services/market/gecko"
	"github.com/everstake/elrond-monitor-backend/smodels"
)

const (
	cmcProvider   = "cmc"
	geckoProvider = "coingecko"
)

type (
	Provider interface {
		GetMarketData() (smodels.MarketData, error)
	}
)

func GetProvider(cfg config.MarketProvider) (p Provider, err error) {
	switch cfg.Title {
	case cmcProvider:
		p = cmc.NewCMC(cfg.APIKey)
	case geckoProvider:
		p = gecko.NewGecko()
	default:
		return nil, errors.New("provider not found")
	}
	return p, nil
}
