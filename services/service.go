package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
)

type (
	Services interface {
		GetTransactions(filter filters.Transactions) (items smodels.Pagination, err error)
		GetTransaction(hash string) (tx smodels.Tx, err error)
		GetBlock(hash string) (block smodels.Block, err error)
		GetBlocks(filter filters.Blocks) (items smodels.Pagination, err error)
		GetBlockByNonce(shard uint64, nonce uint64) (block smodels.Block, err error)
		GetAccounts(filter filters.Accounts) (items smodels.Pagination, err error)
		GetMiniBlock(hash string) (block smodels.Miniblock, err error)
		GetAccount(address string) (account smodels.Account, err error)
		UpdateNodes()
		GetNodes(filter filters.Nodes) (nodes smodels.Pagination, err error)
		UpdateStats()
		GetStats() (stats smodels.Stats, err error)
		GetDailyStats(filter filters.DailyStats) (items []smodels.RangeItem, err error)
		GetEpoch() (epoch smodels.Epoch, err error)
		UpdateValidatorsMap()
		GetValidatorsMap() ([]byte, error)
		GetStakeEvents(filter filters.StakeEvents) (items smodels.Pagination, err error)
		GetStakingProviders(filter filters.StakingProviders) (pagination smodels.Pagination, err error)
		GetStakingProvider(address string) (provider smodels.StakingProvider, err error)
		UpdateStakingProviders()
		GetNode(key string) (node smodels.Node, err error)
		UpdateValidators()
		GetValidators(filter filters.Validators) (pagination smodels.Pagination, err error)
		GetValidator(identity string) (validator smodels.Identity, err error)
		GetValidatorStats() (stats smodels.ValidatorStats, err error)
		MakeRanking()
		GetRanking() (items []smodels.Ranking, err error)
		UpdateTokens()
		GetToken(id string) (token smodels.Token, err error)
		GetTokens(filter filters.Tokens) (pagination smodels.Pagination, err error)
		GetNFTCollection(id string) (collection smodels.NFTCollection, err error)
		GetNFTCollections(filter filters.NFTCollections) (pagination smodels.Pagination, err error)
		GetNFT(id string) (sNFT smodels.NFT, err error)
		GetNFTs(filter filters.NFTTokens) (pagination smodels.Pagination, err error)
		GetOperations(filter filters.Operations) (items smodels.Pagination, err error)
		GetESDTAccounts(filter filters.ESDT) (items smodels.Pagination, err error)
	}
	parser interface {
		GetDelegations(delegator string) map[string]decimal.Decimal
	}

	ServiceFacade struct {
		dao           dao.DAO
		cfg           config.Config
		node          node.APIi
		networkConfig node.NetworkConfig
		parser        parser
	}
)

func NewServices(d dao.DAO, cfg config.Config, p parser) (svc Services, err error) {
	n := node.NewAPI(cfg.Parser.Node, cfg.Contracts)
	nCfg, err := n.GetNetworkConfig()
	if err != nil {
		return nil, fmt.Errorf("GetNetworkConfig: %s", err.Error())
	}
	return &ServiceFacade{
		dao:           d,
		cfg:           cfg,
		node:          n,
		networkConfig: nCfg,
		parser:        p,
	}, nil
}
