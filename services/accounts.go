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
		userStake, err := s.node.GetUserStake(string(a.UserAccount.AddressBytes()))
		if err != nil {
			return items, fmt.Errorf("node.GetUserStake: %s", err.Error())
		}
		balance := decimal.NewFromBigInt(a.UserAccount.GetBalance(), 0)
		accounts[i] = smodels.Account{
			Address:     string(a.UserAccount.AddressBytes()),
			Balance:     node.ValueToEGLD(balance),
			Nonce:       a.UserAccount.GetNonce(),
			Delegated:   node.ValueToEGLD(userStake.ActiveStake),
			Undelegated: node.ValueToEGLD(userStake.UnstakedStake),
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
	acc, err := s.node.GetAddress(address)
	if err != nil {
		return account, fmt.Errorf("node.GetAddress: %s", err.Error())
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
	return smodels.Account{
		Address:          address,
		Balance:          node.ValueToEGLD(acc.Balance),
		Nonce:            acc.Nonce,
		Delegated:        node.ValueToEGLD(userStake.ActiveStake),
		Undelegated:      node.ValueToEGLD(userStake.UnstakedStake),
		RewardsClaimed:   decimal.Zero, // todo
		ClaimableRewards: node.ValueToEGLD(claimableRewards),
		StakingProviders: stakeProviders,
	}, nil
}
