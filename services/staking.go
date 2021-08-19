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
	"sort"
	"time"
)

const (
	stakingProvidersStorageKey = "staking_providers"
)

func (s *ServiceFacade) GetStakeEvents(filter filters.StakeEvents) (page smodels.Pagination, err error) {
	items, err := s.dao.GetStakeEvents(filter)
	if err != nil {
		return page, fmt.Errorf("dao.GetStakeEvents: %s", err.Error())
	}
	total, err := s.dao.GetStakeEventsTotal(filter)
	if err != nil {
		return page, fmt.Errorf("dao.GetStakeEventsTotal: %s", err.Error())
	}
	events := make([]smodels.StakeEvent, len(items))
	for i, item := range items {
		events[i] = smodels.StakeEvent{
			TxHash:    item.TxHash,
			Type:      item.Type,
			Validator: item.Validator,
			Delegator: item.Delegator,
			Epoch:     item.Epoch,
			Amount:    item.Amount,
			CreatedAt: smodels.NewTime(item.CreatedAt),
		}
	}
	return smodels.Pagination{
		Items: events,
		Count: total,
	}, nil
}

func (s *ServiceFacade) GetStakingProviders() (providers []smodels.StakingProvider, err error) {
	err = s.getCache(stakingProvidersStorageKey, &providers)
	if err != nil {
		return nil, fmt.Errorf("getCache: %s", err.Error())
	}
	return providers, nil
}

func (s *ServiceFacade) GetStakingProvider(address string) (provider smodels.StakingProvider, err error) {
	var providers []smodels.StakingProvider
	err = s.getCache(stakingProvidersStorageKey, &providers)
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
		if n.Owner != "" {
			nodesProviders[n.Owner] = append(nodesProviders[n.Owner], n)
		}
	}
	var identities []smodels.Identity
	err = s.getCache(validatorsStorageKey, &identities)
	if err != nil {
		return fmt.Errorf("getCache(identities): %s", err.Error())
	}
	identitiesNamesMap := make(map[string]smodels.Identity)
	for _, identity := range identities {
		identitiesNamesMap[identity.Identity] = identity
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
		var v smodels.StakingProviderValidator
		if identity, ok := identitiesNamesMap[meta.Iidentity]; ok {
			v = smodels.StakingProviderValidator{
				Name:         identity.Name,
				Locked:       identity.Locked,
				StakePercent: identity.StakePercent,
				Nodes:        identity.Validators,
			}
		}
		p := smodels.StakingProvider{
			Provider:         address,
			ServiceFee:       config.ServiceFee.Div(decimal.New(100, 0)),
			DelegationCap:    node.ValueToEGLD(config.DelegationCap),
			APR:              sp.Apr,
			NumUsers:         numUsers,
			CumulatedRewards: node.ValueToEGLD(reward),
			Identity:         meta.Iidentity,
			Featured:         sp.Featured,
			Name:             meta.Name,
			Validator:        v,
		}
		for _, n := range nodesProviders[address] {
			p.NumNodes++
			p.Stake = p.Stake.Add(n.Stake)
			p.TopUp = p.TopUp.Add(n.TopUp)
			p.Locked = p.Locked.Add(n.Locked)
		}
		if p.NumNodes == 0 || p.Stake.Equal(decimal.Zero) {
			continue
		}
		providers = append(providers, p)
	}
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Locked.GreaterThan(providers[j].Locked)
	})
	err = s.setCache(stakingProvidersStorageKey, providers)
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
