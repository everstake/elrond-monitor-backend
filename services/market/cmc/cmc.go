package cmc

import (
	"github.com/everstake/elrond-monitor-backend/smodels"
)

type CMC struct {
	apiKey string
}

func NewCMC(apiKey string) *CMC {
	return &CMC{
		apiKey: apiKey,
	}
}

func (C CMC) GetMarketData() (smodels.MarketData, error) {
	return smodels.MarketData{}, nil
}
