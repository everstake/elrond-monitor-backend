package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
)

func (s *ServiceFacade) GetAccounts(filter filters.Accounts) (items smodels.Pagination, err error) {
	dAccounts, err := s.dao.GetAccounts(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetAccounts: %s", err.Error())
	}
	accounts := make([]smodels.Account, len(dAccounts))
	for i, a := range dAccounts {
		userStake, err := s.node.GetUserStake(a.Address)
		if err != nil {
			return items, fmt.Errorf("node.GetUserStake: %s", err.Error())
		}
		balance, _ := decimal.NewFromString(a.Balance)
		accounts[i] = smodels.Account{
			Address:         a.Address,
			Balance:         node.ValueToEGLD(balance),
			Nonce:           a.Nonce,
			Delegated:       node.ValueToEGLD(userStake.ActiveStake),
			Undelegated:     node.ValueToEGLD(userStake.UnstakedStake),
			IsSmartContract: a.IsSmartContract,
		}
	}
	total, err := s.dao.GetAccountsCount(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetAccountsCount: %s", err.Error())
	}
	return smodels.Pagination{
		Items: accounts,
		Count: total,
	}, nil
}

func (s *ServiceFacade) GetAccount(address string) (account smodels.Account, err error) {
	acc, err := s.dao.GetAccount(address)
	if err != nil {
		return account, fmt.Errorf("dao.GetAccount: %s", err.Error())
	}
	userStake, err := s.node.GetUserStake(address)
	if err != nil {
		return account, fmt.Errorf("node.GetUserStake: %s", err.Error())
	}
	claimableRewards, err := s.node.GetClaimableRewards(address)
	if err != nil {
		return account, fmt.Errorf("node.GetClaimableRewards: %s", err.Error())
	}
	delegations := s.parser.GetDelegations(address)
	var stakeProviders []smodels.AccountStakingProvider
	for validator, stake := range delegations {
		stakeProviders = append(stakeProviders, smodels.AccountStakingProvider{
			Provider: validator,
			Stake:    stake,
		})
	}
	balance, _ := decimal.NewFromString(acc.Balance)
	return smodels.Account{
		Address:          address,
		Balance:          node.ValueToEGLD(balance),
		Nonce:            acc.Nonce,
		Delegated:        node.ValueToEGLD(userStake.ActiveStake),
		Undelegated:      node.ValueToEGLD(userStake.UnstakedStake),
		RewardsClaimed:   decimal.Zero, // todo
		ClaimableRewards: node.ValueToEGLD(claimableRewards),
		StakingProviders: stakeProviders,
		IsSmartContract:  acc.IsSmartContract,
	}, nil
}
