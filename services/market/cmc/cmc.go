package cmc

import "github.com/everstake/elrond-monitor-backend/services/market"

type CMC struct {
	apiKey string
}

func NewCMC(apiKey string) *CMC {
	return &CMC{
		apiKey: apiKey,
	}
}

func (C CMC) GetMarketData() (market.Data, error) {
	return market.Data{}, nil
}
