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
	"strings"
	"time"
)

const (
	stakingProvidersMapStorageKey = "staking_providers"
	nodesStorageKey               = "nodes"
)

func (s *ServiceFacade) UpdateNodes() {
	err := s.updateNodes()
	if err != nil {
		log.Error("updateNodes: %s", err.Error())
	}
}

func (s *ServiceFacade) updateNodes() error {
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
		shard := n.ComputedShardID
		if n.PeerType == smodels.NodeTypeObserver {
			shard = n.ReceivedShardID
		}
		nodesMap[n.PublicKey] = smodels.Node{
			HeartbeatStatus:    n,
			ValidatorStatistic: node.ValidatorStatistic{ShardID: shard},
		}
	}
	for key, stat := range validatorStatistics {
		n := nodesMap[key]
		n.ValidatorStatistic = stat
		nodesMap[key] = n
	}

	for key, n := range nodesMap {
		status := n.ValidatorStatus
		if status == "" {
			status = n.PeerType
		}
		if status == smodels.NodeTypeObserver {
			n.Type = smodels.NodeTypeObserver
		} else {
			n.Type = smodels.NodeTypeValidator
			if strings.Contains(status, smodels.NodeStatusLeaving) {
				status = smodels.NodeStatusLeaving
			}
		}
		n.Status = status
		if n.TotalUpTimeSec == 0 && n.TotalDownTimeSec == 0 {
			n.UpTime = 0
			n.DownTime = 0
			if n.IsActive {
				n.UpTime = 100
				n.DownTime = 100
			}
		} else {
			n.UpTime = float64(n.TotalUpTimeSec*100) / float64(n.TotalUpTimeSec+n.TotalDownTimeSec)
			n.DownTime = 100 - n.UpTime
		}
		nodesMap[key] = n
	}

	// set from queue
	queue, err := s.node.GetQueue()
	if err != nil {
		return fmt.Errorf("node.GetQueue: %s", err.Error())
	}
	for _, item := range queue {
		n, ok := nodesMap[item.BLS]
		if ok {
			n.Type = smodels.NodeTypeValidator
			n.Status = smodels.NodeStatusQueued
			n.Position = item.Position
		} else {
			n = smodels.Node{
				HeartbeatStatus: node.HeartbeatStatus{
					PublicKey: item.BLS,
				},
				Position: item.Position,
				Type:     smodels.NodeTypeValidator,
				Status:   smodels.NodeStatusQueued,
			}
		}
		nodesMap[item.BLS] = n
	}

	// set owners
	for key, n := range nodesMap {
		if n.Type != smodels.NodeTypeValidator {
			continue
		}
		owner, err := s.node.GetOwner(key)
		if err != nil {
			return fmt.Errorf("node.GetOwner: %s", err.Error())
		}
		n.Owner = owner
		nodesMap[key] = n
	}

	var providers []smodels.StakingProvider
	err = s.getCache(stakingProvidersMapStorageKey, &providers)
	if err != nil {
		return fmt.Errorf("getCache(providers): %s", err.Error())
	}
	findProvider := func(owner string) (smodels.StakingProvider, bool) {
		for _, p := range providers {
			if p.Provider == owner {
				return p, true
			}
		}
		return smodels.StakingProvider{}, false
	}

	for key, n := range nodesMap {
		if n.Type != smodels.NodeTypeValidator {
			continue
		}
		p, found := findProvider(n.Owner)
		if found {
			n.Provider = p.Provider
			if p.Identity != "" {
				n.Identity = p.Identity
			}
			nodesMap[key] = n
		}
	}

	// set stakes
	for key, n := range nodesMap {
		if n.Type != smodels.NodeTypeValidator {
			continue
		}
		address := n.Provider
		if address == "" {
			address = n.Owner
		}
		stake, err := s.node.GetTotalStakedTopUpStakedBlsKeys(address)
		if err != nil {
			log.Warn("updateNodes: node.GetTotalStakedTopUpStakedBlsKeys: %s", err.Error())
			continue
		}
		n.Stake = node.ValueToEGLD(stake.Stake)
		n.TopUp = node.ValueToEGLD(stake.TopUp)
		n.Locked = node.ValueToEGLD(stake.Locked)
		nodesMap[key] = n
	}

	var nodes []smodels.Node
	for key, n := range nodesMap {
		n.PublicKey = key
		nodes = append(nodes, n)
	}

	err = s.setCache(nodesStorageKey, nodes)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
}

func (s *ServiceFacade) GetNodes(filter filters.Nodes) (pagination smodels.Pagination, err error) {
	var nodes []smodels.Node
	err = s.getCache(nodesStorageKey, &nodes)
	if err != nil {
		return pagination, fmt.Errorf("getCache: %s", err.Error())
	}
	nodesLen := uint64(len(nodes))
	pagination.Count = nodesLen
	if filter.Limit*(filter.Page-1) > nodesLen {
		return pagination, nil
	}
	maxIndex := filter.Page*filter.Limit - 1
	if nodesLen-1 < maxIndex {
		maxIndex = nodesLen - 1
	}
	pagination.Items = nodes[filter.Offset():maxIndex]
	return pagination, nil
}

func (s *ServiceFacade) GetNode(key string) (node smodels.Node, err error) {
	var nodes []smodels.Node
	err = s.getCache(nodesStorageKey, &nodes)
	if err != nil {
		return node, fmt.Errorf("getCache: %s", err.Error())
	}
	for _, n := range nodes {
		if n.PublicKey == key {
			return n, nil
		}
	}
	msg := fmt.Sprintf("node %s not found", key)
	return node, smodels.Error{
		Err:      msg,
		Msg:      msg,
		HttpCode: 404,
	}
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
			log.Warn("updateStakingProviders: node.GetProviderMeta: %s", err.Error())
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
