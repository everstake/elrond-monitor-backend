package dao

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/postgres"
)

type (
	Postgres interface {
		// parsers
		GetParsers() (parsers []dmodels.Parser, err error)
		GetParser(title string) (parser dmodels.Parser, err error)
		UpdateParser(parser dmodels.Parser) error

		// accounts
		CreateAccounts(accounts []dmodels.Account) error
		GetAccounts() (accounts []dmodels.Account, err error)
		GetAccountsTotal() (total uint64, err error)
		GetAccount(address string) (account dmodels.Account, err error)

		// blocks
		CreateBlocks(blocks []dmodels.Block) error
		CreateMiniBlocks(blocks []dmodels.MiniBlock) error // type == TxBlock

		// transcations
		CreateTransactions(transactions []dmodels.Transaction) error
		CreateSCResults(results []dmodels.SCResult) error
	}

	DAO interface {
		Postgres
	}

	daoImpl struct {
		Postgres
	}
)

func NewDAO(cfg config.Config) (DAO, error) {
	postgresDB, err := postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewPostgres: %s", err.Error())
	}
	return daoImpl{
		Postgres: postgresDB,
	}, nil
}