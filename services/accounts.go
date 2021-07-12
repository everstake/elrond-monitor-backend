package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/smodels"
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
			Balance:         0,
			Delegated:       0,
			Undelegated:     0,
			RewardsClaimed:  0,
			StakingProvider: 0,
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
