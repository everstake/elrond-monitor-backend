package services

import (
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/smodels"
)

type (
	Services interface {
		GetTransactions(filter filters.Transactions) (txs []smodels.Tx, err error)
		GetTransaction(hash string) (tx smodels.Tx, err error)
	}


	ServiceFacade struct {
		dao  dao.DAO
		cfg  config.Config
	}
)

func NewServices(d dao.DAO, cfg config.Config) (svc Services, err error) {
	return &ServiceFacade{
		dao:  d,
		cfg:  cfg,
	}, nil
}

