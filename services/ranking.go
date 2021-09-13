package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
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
func (s *ServiceFacade) makeRanking() error {
	delegations, err := s.dao.GetDelegationState()
	if err != nil {
		return fmt.Errorf("dao.GetDelegationState: %s", err.Error())
	}
	delegationsMap := make(map[string]map[string]decimal.Decimal) // [provider][delegator]amount
	for _, delegation := range delegations {
		if delegation.Amount.IsZero() {
			continue
		}
		if _, ok := delegationsMap[delegation.Validator]; !ok {
			delegationsMap[delegation.Validator] = make(map[string]decimal.Decimal)
		}
		delegationsMap[delegation.Validator][delegation.Delegator] = delegation.Amount
	}

	var providers []smodels.StakingProvider
	err = s.getCache(dmodels.StakingProvidersStorageKey, &providers)
	if err != nil {
		return fmt.Errorf("getCache: %s", err.Error())
	}

	if len(providers) > 100 {
		providers = providers[:100]
	}

	var ranking []smodels.Ranking
	for _, p := range providers {
		rank := smodels.Ranking{Name: p.Name, Address: p.Provider}
		total := decimal.Zero
		for _, amount := range delegationsMap[p.Provider] {
			switch true {
			case amount.GreaterThanOrEqual(decimal.Zero) && amount.LessThan(intToDec(100)):
				rank.T100.Count++
			case amount.GreaterThanOrEqual(intToDec(100)) && amount.LessThan(intToDec(1000)):
				rank.F100T1k.Amount = rank.F100T1k.Amount.Add(amount)
				rank.F100T1k.Count++
			case amount.GreaterThanOrEqual(intToDec(1000)) && amount.LessThan(intToDec(10000)):
				rank.F1kT10k.Amount = rank.F1kT10k.Amount.Add(amount)
				rank.F1kT10k.Count++
			case amount.GreaterThanOrEqual(intToDec(10000)) && amount.LessThan(intToDec(100000)):
				rank.F10kT100k.Amount = rank.F10kT100k.Amount.Add(amount)
				rank.F10kT100k.Count++
			case amount.GreaterThanOrEqual(intToDec(100000)):
				rank.F100k.Amount = rank.F100k.Amount.Add(amount)
				rank.F100k.Count++
			}
			total = total.Add(amount)
		}
		if total.GreaterThan(p.Locked) {
			log.Warn("makeRanking: provider[%s, %s]: total > locked", p.Name, p.Provider)
			continue
		}
		rank.T100.Amount = p.Locked.Sub(total)
		ranking = append(ranking, rank)
	}
	err = s.setCache(dmodels.RankingStorageKey, ranking)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
}

func intToDec(v int64) decimal.Decimal {
	return decimal.New(v, 0)
}
