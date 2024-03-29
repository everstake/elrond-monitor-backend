package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/derrors"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"net/http"
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
			Address:     a.Address,
			Balance:     node.ValueToEGLD(balance),
			Nonce:       a.Nonce,
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
	acc, err := s.dao.GetAccount(address)
	if err != nil {
		if err == derrors.NotFound {
			return account, smodels.Error{
				Err:      err.Error(),
				Msg:      "account not found",
				HttpCode: http.StatusNotFound,
			}
		}
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
	}, nil
}

func (s *ServiceFacade) GetESDTAccounts(filter filters.ESDT) (items smodels.Pagination, err error) {
	accounts, err := s.dao.GetESDTAccounts(filter)
	if err != nil {
		return items, errors.Wrap(err, "get esdt accounts")
	}
	total, err := s.dao.GetESDTAccountsCount(filter)
	if err != nil {
		return items, errors.Wrap(err, "get total esdt accounts")
	}
	var esdtTokens []string
	for _, acc := range accounts {
		found := false
		for _, t := range esdtTokens {
			if t == acc.Token {
				found = true
				break
			}
		}
		if !found {
			esdtTokens = append(esdtTokens, acc.Token)
		}
	}
	esdtTokensMap := make(map[string]dmodels.Token)
	if len(esdtTokens) > 0 {
		tokens, err := s.dao.GetTokens(filters.Tokens{Identifier: esdtTokens})
		if err != nil {
			return items, errors.Wrap(err, "get tokens")
		}
		for _, t := range tokens {
			esdtTokensMap[t.Identity] = t
		}
	}
	acs := make([]smodels.ESDTAccount, len(accounts))
	for i, acc := range accounts {
		mToken := smodels.TokenMetaInfo{
			Identifier: acc.Token,
			Name:       acc.Token,
			Value:      acc.Balance,
		}
		eT, ok := esdtTokensMap[acc.Token]
		if ok {
			mToken.Name = eT.Name
			mToken.Decimal = eT.Decimals
		}
		acs[i] = smodels.ESDTAccount{
			Address: acc.Address,
			Balance: acc.Balance,
			Token:   mToken,
		}
	}
	return smodels.Pagination{
		Items: acs,
		Count: total,
	}, nil
}
