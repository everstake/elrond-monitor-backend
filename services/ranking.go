package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
	"sort"
)

func (s *ServiceFacade) GetRanking() (items []smodels.Ranking, err error) {
	err = s.getCache(dmodels.RankingStorageKey, &items)
	if err != nil {
		return items, fmt.Errorf("getCache: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) MakeRanking() {
	err := s.makeRanking()
	if err != nil {
		log.Error("makeRanking: %s", err.Error())
	}
}

type delegator struct {
	address string
	amount  decimal.Decimal
}

func (s *ServiceFacade) makeRanking() error {
	delegations, err := s.dao.GetDelegationState()
	if err != nil {
		return fmt.Errorf("dao.GetDelegationState: %s", err.Error())
	}
	delegatorsMap := make(map[string]decimal.Decimal)
	delegationsMap := make(map[string]map[string]decimal.Decimal) // [provider][delegator]amount
	for _, delegation := range delegations {
		delegatorsMap[delegation.Delegator] = delegatorsMap[delegation.Delegator].Add(delegation.Amount)
		if _, ok := delegationsMap[delegation.Validator]; !ok {
			delegationsMap[delegation.Validator] = make(map[string]decimal.Decimal)
		}
		delegationsMap[delegation.Validator][delegation.Delegator] = delegation.Amount
	}
	delegators := make([]delegator, 0, len(delegatorsMap))
	for address, amount := range delegatorsMap {
		delegators = append(delegators, delegator{
			address: address,
			amount:  amount,
		})
	}
	sort.Slice(delegators, func(i, j int) bool {
		return delegators[i].amount.GreaterThan(delegations[j].Amount)
	})
	if len(delegators) > 200 {
		delegators = delegators[:200]
	}
	addressesMap := make(map[string]bool)
	for _, d := range delegators {
		addressesMap[d.address] = true
	}
	var providers []smodels.StakingProvider
	err = s.setCache(dmodels.StakingProvidersStorageKey, providers)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	if len(providers) > 100 {
		providers = providers[:100]
	}
	var ranking []smodels.Ranking
	for _, provider := range providers {
		name := provider.Name
		if name == "" {
			name = provider.Identity
		}
		if name == "" {
			name = provider.Provider
		}
		ranking = append(ranking, smodels.Ranking{
			Provider:   name,
			Amount:     provider.Locked,
			Delegators: delegationsMap[provider.Provider],
		})
	}
	err = s.setCache(dmodels.RankingStorageKey, ranking)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
}
