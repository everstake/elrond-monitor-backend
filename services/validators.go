package services

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

const (
	stakingProvidersMapStorageKey = "staking_providers"
	nodesStorageKey               = "nodes"
)

func (s *ServiceFacade) UpdateNodes() error {
	nodesStatus, err := s.node.GetHeartbeatStatus()
	if err != nil {
		return fmt.Errorf("node.GetHeartbeatStatus: %s", err.Error())
	}
	validatorStatistics, err := s.node.GetValidatorStatistics()
	if err != nil {
		return fmt.Errorf("node.GetValidatorStatistics: %s", err.Error())
	}
	nodesMap := make(map[string]smodels.Node)
	for _, n := range nodesStatus {
		t := smodels.NodeTypeObserver
		if n.PeerType != smodels.NodeTypeObserver {
			t = smodels.NodeTypeValidator
		}
		nodesMap[n.PublicKey] = smodels.Node{HeartbeatStatus: n, Type: t}

	}
	for key, stat := range validatorStatistics {
		n := nodesMap[key]
		n.ValidatorStatistic = stat
		nodesMap[key] = smodels.Node{ValidatorStatistic: stat}
	}

	for key, node := range nodesMap {
		if node.TotalUpTimeSec == 0 && node.TotalDownTimeSec == 0 {
			node.UpTime = 0
			node.DownTime = 0
			if node.IsActive {
				node.UpTime = 100
				node.DownTime = 100
			}
		} else {
			node.UpTime = float64(node.TotalUpTimeSec*100) / float64(node.TotalUpTimeSec+node.TotalDownTimeSec)
			node.DownTime = 100 - node.UpTime
		}
		nodesMap[key] = node
	}

	//providers, err := s.GetStakingProviders()
	//if err != nil {
	//	return fmt.Errorf("GetStakingProviders: %s", err.Error())
	//}

	//for _, p := range providers {
	//	node, ok := nodesMap[p. ?]
	//	if ok && node.Type == smodels.NodeTypeValidator {
	//		node.Owner = p.Owner
	//		node.Provider = p.Contract
	//		if p.Identity.Name != "" {
	//
	//		}
	//	}
	//}

	err = s.setCache(nodesStorageKey, nodesStatus)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
}

func (s *ServiceFacade) GetNodes(filter filters.Nodes) (nodes []smodels.Node, err error) {
	err = s.getCache(nodesStorageKey, &nodes)
	if err != nil {
		return nil, fmt.Errorf("getCache: %s", err.Error())
	}
	nodesLen := uint64(len(nodes))
	if filter.Limit*(filter.Page-1) > nodesLen {
		return nil, nil
	}
	maxIndex := filter.Page*filter.Limit - 1
	if nodesLen-1 < maxIndex {
		maxIndex = nodesLen - 1
	}
	return nodes, nil
}

func (s *ServiceFacade) GetStakingProviders() (providers []smodels.StakingProvider, err error) {
	err = s.getCache(stakingProvidersMapStorageKey, &providers)
	if err != nil {
		return nil, fmt.Errorf("getCache: %s", err.Error())
	}
	return providers, nil
}

func (s *ServiceFacade) GetStakingProvider(address string) (provider smodels.StakingProvider, err error) {
	var providers []smodels.StakingProvider
	err = s.getCache(stakingProvidersMapStorageKey, &providers)
	if err != nil {
		return provider, fmt.Errorf("getCache: %s", err.Error())
	}
	for _, p := range providers {
		if p.Provider == address {
			return p, err
		}
	}
	msg := fmt.Sprintf("provider %s not found", address)
	return provider, smodels.Error{
		Err:      msg,
		Msg:      msg,
		HttpCode: 404,
	}
}

func (s *ServiceFacade) UpdateStakingProviders() {
	err := s.updateStakingProviders()
	if err != nil {
		log.Error("updateStakingProviders: %s", err.Error())
	}
}

func (s *ServiceFacade) updateStakingProviders() error {
	addresses, err := s.node.GetProviderAddresses()
	if err != nil {
		return fmt.Errorf("node.GetProviderAddresses: %s", err.Error())
	}
	sourceProviders, err := s.getStakingProvidersFromSource()
	if err != nil {
		log.Error("updateStakingProviders: getStakingProvidersFromSource: %s", err.Error())
	}
	sourceProvidersMap := make(map[string]smodels.SourceStakingProvider)
	for _, provider := range sourceProviders {
		sourceProvidersMap[provider.Contract] = provider
	}
	providers := make([]smodels.StakingProvider, len(addresses))
	for i, address := range addresses {
		config, err := s.node.GetProviderConfig(address)
		if err != nil {
			return fmt.Errorf("node.GetProviderConfig: %s", err.Error())
		}
		meta, err := s.node.GetProviderMeta(address)
		if err != nil {
			return fmt.Errorf("node.GetProviderMeta: %s", err.Error())
		}
		numUsers, err := s.node.GetProviderNumUsers(address)
		if err != nil {
			return fmt.Errorf("node.GetProviderNumUsers: %s", err.Error())
		}
		reward, err := s.node.GetCumulatedRewards(address)
		if err != nil {
			return fmt.Errorf("node.GetCumulatedRewards: %s", err.Error())
		}
		sp := sourceProvidersMap[address]
		providers[i] = smodels.StakingProvider{
			Provider:         address,
			ServiceFee:       config.ServiceFee.Div(decimal.New(100, 0)),
			DelegationCap:    node.ValueToEGLD(config.DelegationCap),
			APR:              sp.Apr,
			NumUsers:         numUsers,
			CumulatedRewards: node.ValueToEGLD(reward),
			Identity:         meta.Iidentity,
			NumNodes:         0,
			Stake:            decimal.Decimal{},
			TopUp:            decimal.Decimal{},
			Locked:           decimal.Decimal{},
			Featured:         sp.Featured,
		}
	}
	err = s.setCache(stakingProvidersMapStorageKey, providers)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
}

func (s *ServiceFacade) getStakingProvidersFromSource() (providers []smodels.SourceStakingProvider, err error) {
	client := http.Client{Timeout: time.Second * 30}
	resp, err := client.Get(s.cfg.StakingProvidersSource)
	if err != nil {
		return providers, fmt.Errorf("client.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return providers, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&providers)
	if err != nil {
		return providers, fmt.Errorf("json.Decode: %s", err.Error())
	}
	return providers, nil
}
