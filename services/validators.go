package services

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"net/http"
	"time"
)

const stakingProvidersMapStorageKey = "staking_providers"

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
		if p.Contract == address {
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
	client := http.Client{Timeout: time.Second * 30}
	resp, err := client.Get(s.cfg.StakingProvidersSource)
	if err != nil {
		return fmt.Errorf("client.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	var providers []smodels.StakingProvider
	err = json.NewDecoder(resp.Body).Decode(&providers)
	if err != nil {
		return fmt.Errorf("json.Decode: %s", err.Error())
	}
	err = s.setCache(stakingProvidersMapStorageKey, providers)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
}
