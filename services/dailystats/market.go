package dailystats

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func (ds *DailyStats) GetMarket() (map[string]decimal.Decimal, error) {
	data, err := ds.market.GetMarketData()
	if err != nil {
		return nil, fmt.Errorf("market.GetMarketData: %s", err.Error())
	}
	return map[string]decimal.Decimal{
		PriceKey:         data.Price,
		TradingVolumeKey: data.TradingVolume24h,
	}, nil
}
