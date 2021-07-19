package services

import (
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
)

type (
	Services interface {
		GetTransactions(filter filters.Transactions) (items smodels.Pagination, err error)
		GetTransaction(hash string) (tx smodels.Tx, err error)
		GetBlock(hash string) (block smodels.Block, err error)
		GetBlocks(filter filters.Blocks) (items smodels.Pagination, err error)
		GetAccounts(filter filters.Accounts) (items smodels.Pagination, err error)
		GetMiniBlock(hash string) (block smodels.Miniblock, err error)
		GetAccount(address string) (account smodels.Account, err error)
	}

	ServiceFacade struct {
		dao  dao.DAO
		cfg  config.Config
		node node.APIi
	}
)

func NewServices(d dao.DAO, cfg config.Config) (svc Services, err error) {
	return &ServiceFacade{
		dao:  d,
		cfg:  cfg,
		node: node.NewAPI(cfg.Parser.Node),
	}, nil
}
