package cmc

import (
	smodels2 "github.com/everstake/elrond-monitor-backend/smodels"
)

type CMC struct {
	apiKey string
}

func NewCMC(apiKey string) *CMC {
	return &CMC{
		apiKey: apiKey,
	}
}

func (C CMC) GetMarketData() (smodels2.MarketData, error) {
	return smodels2.MarketData{}, nil
}
