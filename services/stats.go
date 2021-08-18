package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services/market"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"io/ioutil"
	"net/http"
)

const (
	statsStorageKey         = "stats"
	validatorsMapSource     = "https://internal-api.elrond.com/markers"
	validatorsMapStorageKey = "validators_map"
)

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
	accountsTotal, err := s.dao.GetAccountsCount(filters.Accounts{})
	if err != nil {
		return fmt.Errorf("dao.GetAccountsCount: %s", err.Error())
	}
	txsTotal, err := s.dao.GetTransactionsCount(filters.Transactions{})
	if err != nil {
		return fmt.Errorf("dao.GetTransactionsCount: %s", err.Error())
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

func (s *ServiceFacade) GetValidatorsMap() ([]byte, error) {
	data, err := s.dao.GetStorageValue(validatorsMapStorageKey)
	if err != nil {
		return nil, fmt.Errorf("dao.GetStorageValue: %s", err.Error())
	}
	return []byte(data), nil
}

func (s *ServiceFacade) UpdateValidatorsMap() {
	err := s.updateValidatorsMap()
	if err != nil {
		log.Error("updateValidatorsMap: %s", err.Error())
	}
}

func (s *ServiceFacade) updateValidatorsMap() error {
	resp, err := http.DefaultClient.Get(validatorsMapSource)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Get: %s", err.Error())
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	err = s.dao.UpdateStorageValue(dmodels.StorageItem{
		Key:   validatorsMapStorageKey,
		Value: string(data),
	})
	if err != nil {
		return fmt.Errorf("dao.UpdateStorageValue: %s", err.Error())
	}
	return nil
}
