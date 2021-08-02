package gecko

import (
	"errors"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
	coingecko "github.com/superoo7/go-gecko/v3"
	"net/http"
	"time"
)

const (
	elrondCoinID  = "elrond-erd-2"
	quoteCurrency = "usd"
)

type Gecko struct {
	client *coingecko.Client
}

func NewGecko() *Gecko {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return &Gecko{
		client: coingecko.NewClient(httpClient),
	}
}

func (g Gecko) GetMarketData() (d smodels.MarketData, err error) {
	data, err := g.client.CoinsID(elrondCoinID, false, true, true, false, false, false)
	if err != nil {
		return d, fmt.Errorf("client.CoinsID: %s", err.Error())
	}
	if data.MarketData.MarketCap == nil {
		return d, errors.New("MarketData.MarketCap is nil")
	}
	totalSupply := decimal.Zero
	if data.MarketData.TotalSupply != nil {
		totalSupply = decimal.NewFromFloat(*data.MarketData.TotalSupply)
	}
	return smodels.MarketData{
		Price:             decimal.NewFromFloat(data.MarketData.CurrentPrice[quoteCurrency]),
		PriceChange:       decimal.NewFromFloat(data.MarketData.PriceChange24h),
		Cap:               decimal.NewFromFloat(data.MarketData.MarketCap[quoteCurrency]),
		CapChange:         decimal.NewFromFloat(data.MarketData.MarketCapChangePercentage24h),
		TradingVolume24h:  decimal.NewFromFloat(data.MarketData.TotalVolume[quoteCurrency]),
		CirculatingSupply: decimal.NewFromFloat(data.MarketData.CirculatingSupply),
		TotalSupply:       totalSupply,
	}, nil
}
