package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
)

func (db Postgres) CreateAccounts(accounts []dmodels.Account) error {
	if len(accounts) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.AccountsTable).Columns(
		"acc_address",
		"acc_created_at",
	)
	for _, account := range accounts {
		if account.Address == "" {
			return fmt.Errorf("field Address is empty")
		}
		if account.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			account.Address,
			account.CreatedAt,
		)
	}
	_, err := db.insert(q)
	return err
}

func (db Postgres) GetAccounts() (accounts []dmodels.Account, err error) {
	q := squirrel.Select("*").From(dmodels.AccountsTable)
	err = db.find(&accounts, q)
	return accounts, err
}

func (db Postgres) GetAccountsTotal() (total uint64, err error) {
	q := squirrel.Select("count(*) as total").From(dmodels.AccountsTable)
	err = db.first(&total, q)
	return total, err
}

func (db Postgres) GetAccount(address string) (account dmodels.Account, err error) {
	q := squirrel.Select("*").From(dmodels.AccountsTable).Where(squirrel.Eq{"acc_address": address})
	err = db.first(&account, q)
	return account, err
}
