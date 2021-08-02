package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services/market"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
)

const statsStorageKey = "stats"

func (s *ServiceFacade) GetStats() (stats smodels.Stats, err error) {
	err = s.getCache(statsStorageKey, &stats)
	if err != nil {
		return stats, fmt.Errorf("getCache: %s", err.Error())
	}
	return stats, nil
}

func (s *ServiceFacade) UpdateStats() {
	err := s.updateStats()
	if err != nil {
		log.Error("updateStats: %s", err.Error())
	}
}

func (s *ServiceFacade) updateStats() error {
	m, err := market.GetProvider(s.cfg.MarketProvider)
	if err != nil {
		return fmt.Errorf("market.GetProvider: %s", err.Error())
	}
	marketData, err := m.GetMarketData()
	if err != nil {
		return fmt.Errorf("market.GetMarketData: %s", err.Error())
	}
	status, err := s.node.GetNetworkStatus(node.MetaChainShardIndex)
	if err != nil {
		return fmt.Errorf("node.GetNetworkStatus: %s", err.Error())
	}
	accountsTotal, err := s.dao.GetAccountsTotal(filters.Accounts{})
	if err != nil {
		return fmt.Errorf("dao.GetAccountsTotal: %s", err.Error())
	}
	txsTotal, err := s.dao.GetTransactionsTotal(filters.Transactions{})
	if err != nil {
		return fmt.Errorf("dao.GetTransactionsTotal: %s", err.Error())
	}
	err = s.setCache(statsStorageKey, smodels.Stats{
		Price:             marketData.Price,
		PriceChange:       marketData.PriceChange,
		TradingVolume:     marketData.TradingVolume24h,
		Cap:               marketData.Cap,
		CapChange:         marketData.CapChange,
		CirculatingSupply: marketData.CirculatingSupply,
		TotalSupply:       marketData.TotalSupply,
		Height:            status.ErdNonce,
		TotalTxs:          txsTotal,
		TotalAccounts:     accountsTotal,
	})
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
}
