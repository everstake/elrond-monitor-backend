package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
	"net/http"
	"sort"
	"time"
)

func (s *ServiceFacade) GetValidators(filter filters.Validators) (pagination smodels.Pagination, err error) {
	var validators []smodels.Identity
	err = s.getCache(dmodels.ValidatorsStorageKey, &validators)
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
	err = s.getCache(dmodels.ValidatorsStorageKey, &validators)
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
	err := s.getCache(dmodels.NodesStorageKey, &nodes)
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
		providersMap := make(map[string]interface{})
		for _, n := range ns {
			stake = stake.Add(n.Stake)
			topUp = topUp.Add(n.TopUp)
			score += n.RatingModifier
			if n.Type == smodels.NodeTypeValidator && n.Status != "inactive" {
				count++
			}
			if n.Provider != "" {
				providersMap[n.Provider] = nil
			}
		}
		var providers []string
		for p := range providersMap {
			providers = append(providers, p)
		}
		locked := stake.Add(topUp)
		stakePercent, _ := locked.Mul(decimal.New(100, 0)).Div(totalLocked).Float64()
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
	err = s.setCache(dmodels.ValidatorsStorageKey, identities)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
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
