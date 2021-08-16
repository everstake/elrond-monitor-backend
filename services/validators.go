package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	stakingProvidersMapStorageKey = "staking_providers"
	nodesStorageKey               = "nodes"
	validatorsStorageKey          = "validators"
)

func (s *ServiceFacade) GetValidators(filter filters.Validators) (pagination smodels.Pagination, err error) {
	var validators []smodels.Identity
	err = s.getCache(validatorsStorageKey, &validators)
	if err != nil {
		return pagination, fmt.Errorf("getCache: %s", err.Error())
	}
	validatorsLen := uint64(len(validators))
	pagination.Count = validatorsLen
	if filter.Limit*(filter.Page-1) > validatorsLen {
		return pagination, nil
	}
	maxIndex := filter.Page * filter.Limit
	if validatorsLen-1 < maxIndex {
		maxIndex = validatorsLen - 1
	}
	pagination.Items = validators[filter.Offset():maxIndex]
	return pagination, nil
}

func (s *ServiceFacade) GetValidator(identity string) (validator smodels.Identity, err error) {
	var validators []smodels.Identity
	err = s.getCache(validatorsStorageKey, &validators)
	if err != nil {
		return validator, fmt.Errorf("getCache: %s", err.Error())
	}
	for _, n := range validators {
		if n.Identity == identity {
			return n, nil
		}
	}
	msg := fmt.Sprintf("validator %s not found", identity)
	return validator, smodels.Error{
		Err:      msg,
		Msg:      msg,
		HttpCode: 404,
	}
}

func (s *ServiceFacade) UpdateValidators() {
	err := s.updateValidators()
	if err != nil {
		log.Error("updateValidators: %s", err.Error())
	}
}

func (s *ServiceFacade) updateValidators() error {
	var nodes []smodels.Node
	err := s.getCache(nodesStorageKey, &nodes)
	if err != nil {
		return fmt.Errorf("getCache(nodes): %s", err.Error())
	}
	identitiesMap := make(map[string][]smodels.Node)
	var totalStake decimal.Decimal
	var totalTopUp decimal.Decimal
	for _, n := range nodes {
		if n.Identity != "" {
			identitiesMap[n.Identity] = append(identitiesMap[n.Identity], n)
		} else {
			identitiesMap[n.PublicKey] = []smodels.Node{n}
		}
		if n.Type == smodels.NodeTypeValidator {
			totalStake = totalStake.Add(n.Stake)
			totalTopUp = totalTopUp.Add(n.TopUp)
		}
	}
	totalLocked := totalStake.Add(totalTopUp)
	if totalLocked.IsZero() {
		return errors.New("updateValidators: totalLocked is zero")
	}
	var identities []smodels.Identity
	for key, ns := range identitiesMap {
		var stake, topUp decimal.Decimal
		var score float64
		var count uint64
		var providers []string
		for _, n := range ns {
			stake = stake.Add(n.Stake)
			topUp = topUp.Add(n.TopUp)
			score += n.RatingModifier
			if n.Type == smodels.NodeTypeValidator && n.Status != "inactive" {
				count++
			}
			if n.Provider != "" {
				providers = append(providers, n.Provider)
			}
		}
		locked := stake.Add(topUp)
		stakePercent, _ := locked.Mul(decimal.New(10000, 0)).Div(totalLocked).Float64()
		var kb smodels.IdentityKeybase
		if len(key) < 192 && len(key) != 0 {
			kb, err = s.getIdentityProfile(key)
			if err != nil {
				log.Warn("updateValidators: getIdentityProfile(%s): %s", key, err.Error())
			}
		}
		identities = append(identities, smodels.Identity{
			Avatar:       kb.Them.Pictures.Primary.URL,
			Description:  kb.Them.Profile.Bio,
			Identity:     key,
			Locked:       locked,
			Name:         kb.Them.Profile.FullName,
			Score:        uint64(score),
			Stake:        stake,
			StakePercent: stakePercent,
			TopUp:        topUp,
			Validators:   count,
			Providers:    providers,
		})
	}
	sort.Slice(identities, func(i, j int) bool {
		return identities[i].StakePercent > identities[j].StakePercent
	})
	for i := range identities {
		identities[i].Rank = uint64(i + 1)
	}
	err = s.setCache(validatorsStorageKey, identities)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
}

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
		log.Warn("updateNodes: getCache(providers): %s", err.Error())
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
	maxIndex := filter.Page * filter.Limit
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
	var nodes []smodels.Node
	err = s.getCache(nodesStorageKey, &nodes)
	if err != nil {
		return fmt.Errorf("getCache(nodes): %s", err.Error())
	}
	nodesProviders := make(map[string][]smodels.Node)
	for _, n := range nodes {
		if n.Provider != "" {
			nodesProviders[n.Provider] = append(nodesProviders[n.Provider], n)
		}
	}
	var providers []smodels.StakingProvider
	for _, address := range addresses {
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
		p := smodels.StakingProvider{
			Provider:         address,
			ServiceFee:       config.ServiceFee.Div(decimal.New(100, 0)),
			DelegationCap:    node.ValueToEGLD(config.DelegationCap),
			APR:              sp.Apr,
			NumUsers:         numUsers,
			CumulatedRewards: node.ValueToEGLD(reward),
			Identity:         meta.Iidentity,
			Featured:         sp.Featured,
		}
		for _, n := range nodesProviders[address] {
			p.NumNodes++
			p.Stake = p.Stake.Add(n.Stake)
			p.TopUp = p.Stake.Add(n.TopUp)
			p.Locked = p.Stake.Add(n.Locked)
		}
		if p.NumNodes == 0 || p.Stake.Equal(decimal.Zero) {
			continue
		}
		providers = append(providers, p)
	}
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Locked.GreaterThan(providers[j].Locked)
	})
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

func (s *ServiceFacade) getIdentityProfile(identity string) (data smodels.IdentityKeybase, err error) {
	client := http.Client{Timeout: time.Second * 30}
	resp, err := client.Get(fmt.Sprintf("https://keybase.io/_/api/1.0/user/lookup.json?username=%s", identity))
	if err != nil {
		return data, fmt.Errorf("client.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, fmt.Errorf("json.Decode: %s", err.Error())
	}
	return data, nil
}
