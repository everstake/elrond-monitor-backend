package dao

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/dao/postgres"
)

type (
	Postgres interface {
		// parsers
		GetParsers() (parsers []dmodels.Parser, err error)
		GetParser(title string) (parser dmodels.Parser, err error)
		UpdateParserHeight(parser dmodels.Parser) error

		// accounts
		CreateAccounts(accounts []dmodels.Account) error
		// todo: GetAccounts add pagination
		GetAccounts() (accounts []dmodels.Account, err error)
		GetAccountsTotal() (total uint64, err error)
		GetAccount(address string) (account dmodels.Account, err error)

		// blocks
		CreateBlocks(blocks []dmodels.Block) error
		CreateMiniBlocks(blocks []dmodels.MiniBlock) error // type == TxBlock
		GetBlocks(filter filters.Blocks) (blocks []dmodels.Block, err error)
		GetBlock(hash string) (block dmodels.Block, err error)
		GetMiniBlocks(filter filters.MiniBlocks) (blocks []dmodels.MiniBlock, err error)

		// transcations
		CreateTransactions(transactions []dmodels.Transaction) error
		CreateSCResults(results []dmodels.SCResult) error
		GetTransactions(filter filters.Transactions) (txs []dmodels.Transaction, err error)
		GetTransaction(hash string) (tx dmodels.Transaction, err error)
		GetSCResults(txHash string) (results []dmodels.SCResult, err error)
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
