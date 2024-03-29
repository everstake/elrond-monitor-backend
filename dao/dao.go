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

		// storage
		GetStorageValue(key string) (value string, err error)
		UpdateStorageValue(item dmodels.StorageItem) error

		// staking
		CreateDelegations(delegations []dmodels.Delegation) error

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

		// tokens
		CreateToken(token dmodels.Token) error
		UpdateToken(token dmodels.Token) error
		GetTokens(filter filters.Tokens) (tokens []dmodels.Token, err error)
		GetTokensCount(filter filters.Tokens) (total uint64, err error)
		GetToken(ident string) (token dmodels.Token, err error)

		// nft collections
		CreateNFTCollection(collection dmodels.NFTCollection) error
		UpdateNFTCollection(collection dmodels.NFTCollection) error
		GetNFTCollections(filter filters.NFTCollections) (collections []dmodels.NFTCollection, err error)
		GetNFTCollectionsTotal(filter filters.NFTCollections) (total uint64, err error)
		GetNFTCollection(ident string) (collection dmodels.NFTCollection, err error)
	}

	ElasticSearch interface {
		GetBlock(hash string) (block data.Block, err error)
		GetTransaction(hash string) (tx es.Tx, err error)
		GetMiniblock(hash string) (miniblock data.Miniblock, err error)
		GetBlocks(filter filters.Blocks) (blocks []data.Block, err error)
		GetBlocksCount(filter filters.Blocks) (total uint64, err error)
		GetTransactions(filter filters.Transactions) (txs []data.Transaction, err error)
		GetSCResults(txHash string) (results []es.SCResult, err error)
		GetTransactionsCount(filter filters.Transactions) (total uint64, err error)
		ValidatorsKeys(shard uint64, epoch uint64) (keys data.ValidatorsPublicKeys, err error)
		GetAccount(address string) (acc data.AccountInfo, err error)
		GetAccounts(filter filters.Accounts) (accounts []data.AccountInfo, err error)
		GetAccountsCount(filter filters.Accounts) (total uint64, err error)
		GetESDTAccounts(filter filters.ESDT) (accounts []es.AccountESDT, err error)
		GetESDTAccountsCount(filter filters.ESDT) (total uint64, err error)
		GetOperations(filter filters.Operations) (txs []es.Operation, err error)
		GetOperationsCount(filter filters.Operations) (total uint64, err error)
		GetTokenInfo(id string) (token data.TokenInfo, err error)
		GetNFTTokens(filter filters.NFTTokens) (txs []data.TokenInfo, err error)
		GetNFTTokensCount(filter filters.NFTTokens) (total uint64, err error)
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
