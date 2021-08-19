package dailystats

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
)

func (ds *DailyStats) GetTotalAccounts() (map[string]decimal.Decimal, error) {
	total, err := ds.dao.GetAccountsCount(filters.Accounts{})
	if err != nil {
		return nil, fmt.Errorf("dao.GetAccountsCount: %s", err.Error())
	}
	return map[string]decimal.Decimal{
		TotalAccountKey: decimal.NewFromInt(int64(total)),
	}, nil
}

func (ds *DailyStats) GetTotalTransactions() (map[string]decimal.Decimal, error) {
	total, err := ds.dao.GetTransactionsCount(filters.Transactions{})
	if err != nil {
		return nil, fmt.Errorf("dao.GetTransactionsCount: %s", err.Error())
	}
	return map[string]decimal.Decimal{
		TotalTransactionsKey: decimal.NewFromInt(int64(total)),
	}, nil
}

func (ds *DailyStats) GetTotalDelegators() (map[string]decimal.Decimal, error) {
	var providers []smodels.StakingProvider
	value, err := ds.dao.GetStorageValue(dmodels.StakingProvidersStorageKey)
	if err != nil {
		return nil, fmt.Errorf("dao.GetStorageValue: %s", err.Error())
	}
	err = json.Unmarshal([]byte(value), &providers)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	var total uint64
	for _, p := range providers {
		total += p.NumUsers
	}
	return map[string]decimal.Decimal{
		TotalDelegatorsKey: decimal.NewFromInt(int64(total)),
	}, nil
}
