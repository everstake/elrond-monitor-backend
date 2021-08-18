package dao

import (
	"fmt"
	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/dao/postgres"
	"github.com/everstake/elrond-monitor-backend/services/es"
)

type (
	Postgres interface {
		// parsers
		GetParsers() (parsers []dmodels.Parser, err error)
		GetParser(title string) (parser dmodels.Parser, err error)
		UpdateParserHeight(parser dmodels.Parser) error

		// accounts
		CreateAccounts(accounts []dmodels.Account) error
		//GetAccounts(filter filters.Accounts) (accounts []dmodels.Account, err error)
		//GetAccountsTotal(filter filters.Accounts) (total uint64, err error)
		//GetAccount(address string) (account dmodels.Account, err error)

		// blocks
		CreateBlocks(blocks []dmodels.Block) error
		CreateMiniBlocks(blocks []dmodels.MiniBlock) error // type == TxBlock
		//GetBlocks(filter filters.Blocks) (blocks []dmodels.Block, err error)
		//GetBlock(hash string) (block dmodels.Block, err error)
		//GetMiniBlocks(filter filters.MiniBlocks) (blocks []dmodels.MiniBlock, err error)
		//GetBlocksTotal(filter filters.Blocks) (total uint64, err error)
		//GetMiniBlocksTotal(filter filters.MiniBlocks) (total uint64, err error)
		//GetMiniBlock(hash string) (block dmodels.MiniBlock, err error)

		// transcations
		CreateTransactions(transactions []dmodels.Transaction) error
		CreateSCResults(results []dmodels.SCResult) error
		//GetTransactions(filter filters.Transactions) (txs []dmodels.Transaction, err error)
		//GetTransaction(hash string) (tx dmodels.Transaction, err error)
		//GetTransactionsTotal(filter filters.Transactions) (total uint64, err error)
		//GetSCResults(txHash string) (results []dmodels.SCResult, err error)

		// storage
		GetStorageValue(key string) (value string, err error)
		UpdateStorageValue(item dmodels.StorageItem) error

		// staking
		CreateDelegations(delegations []dmodels.Delegation) error
		CreateStakes(stakes []dmodels.Stake) error

		// rewards
		CreateRewards(rewards []dmodels.Reward) error

		// stake events
		CreateStakeEvents(events []dmodels.StakeEvent) error
		GetDelegationState() (items []dmodels.StakeState, err error)
		GetStakeState() (items []dmodels.StakeState, err error)
		GetStakeEvents(filter filters.StakeEvents) (items []dmodels.StakeEvent, err error)
		GetStakeEventsTotal(filter filters.StakeEvents) (total uint64, err error)

		// daily stats
		CreateDailyStats(stats []dmodels.DailyStat) error
		GetDailyStatsRange(filter filters.DailyStats) (items []dmodels.DailyStat, err error)
	}

	ElasticSearch interface {
		GetBlock(hash string) (block data.Block, err error)
		GetTransaction(hash string) (tx data.Transaction, err error)
		GetMiniblock(hash string) (miniblock data.Miniblock, err error)
		GetBlocks(filter filters.Blocks) (blocks []data.Block, err error)
		GetBlocksCount(filter filters.Blocks) (total uint64, err error)
		GetTransactions(filter filters.Transactions) (txs []data.Transaction, err error)
		GetTransactionsCount(filter filters.Transactions) (total uint64, err error)
		GetMiniblocks(filter filters.MiniBlocks) (txs []data.Miniblock, err error)
		GetSCResults(txHash string) (scs []data.ScResult, err error)
		ValidatorsKeys(shard uint64, epoch uint64) (keys data.ValidatorsPublicKeys, err error)
		GetAccount(address string) (acc data.AccountInfo, err error)
		GetAccounts(filter filters.Accounts) (accounts []data.AccountInfo, err error)
		GetAccountsCount(filter filters.Accounts) (total uint64, err error)
	}

	DAO interface {
		Postgres
		ElasticSearch
	}

	daoImpl struct {
		Postgres
		ElasticSearch
	}
)

func NewDAO(cfg config.Config) (DAO, error) {
	postgresDB, err := postgres.NewPostgres(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewPostgres: %s", err.Error())
	}
	elastic, err := es.NewClient(cfg.ElasticSearch.Address)
	if err != nil {
		return nil, fmt.Errorf("es.NewClient: %s", err.Error())
	}
	return daoImpl{
		Postgres:      postgresDB,
		ElasticSearch: elastic,
	}, nil
}
