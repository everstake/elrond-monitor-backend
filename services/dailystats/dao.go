package dailystats

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
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
