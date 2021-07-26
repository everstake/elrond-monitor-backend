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
		accounts[i] = smodels.Account{
			Address: a.Address,
			// todo
			Balance:         decimal.Decimal{},
			Delegated:       decimal.Decimal{},
			Undelegated:     decimal.Decimal{},
			RewardsClaimed:  decimal.Decimal{},
			StakingProvider: nil,
			CreatedAt:       smodels.NewTime(a.CreatedAt),
		}
	}
	total, err := s.dao.GetAccountsTotal(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetAccountsTotal: %s", err.Error())
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
	dAcc, _ := s.dao.GetAccount(address)
	//delegation, err := s.node.GetAccountDelegation(address)
	//if err != nil {
	//	return account, fmt.Errorf("node.GetAccountDelegation: %s", err.Error())
	//}
	return smodels.Account{
		Address: address,
		Balance: node.ValueToEGLD(acc.Balance),
		Nonce:   acc.Nonce,
		//Delegated:        node.ValueToEGLD(delegation.UserActiveStake),
		//Undelegated:      node.ValueToEGLD(delegation.UserUnstakedStake),
		//ClaimableRewards: node.ValueToEGLD(delegation.ClaimableRewards),
		RewardsClaimed:  decimal.Decimal{},
		StakingProvider: nil,
		CreatedAt:       smodels.NewTime(dAcc.CreatedAt),
	}, nil
}
