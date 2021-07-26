package dailystats

import (
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services/market"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/shopspring/decimal"
	"reflect"
	"runtime"
	"time"
)

const (
	PriceKey         = "price"
	TradingVolumeKey = "trading_volume"
	TotalStakeKey    = "total_stake"
	TotalFeeKey      = "total_fee"
	TotalSupplyKey   = "total_supply"
	TotalAccountKey  = "total_accounts"
	TopUpAmountKey   = "top_up"
)

type (
	DailyStats struct {
		dao     dao.DAO
		node    node.APIi
		market  market.Provider
		stopSig chan struct{}
		actions []action
	}

	action func() (map[string]decimal.Decimal, error)
)

func NewDailyStats(cfg config.Config, d dao.DAO) *DailyStats {
	ds := &DailyStats{
		dao:     d,
		node:    node.NewAPI(cfg.Parser.Node),
		stopSig: make(chan struct{}),
		market:  market.GetProvider(cfg.MarketProvider),
	}
	ds.actions = []action{
		ds.GetMarket,
	}
	return ds
}

func (ds *DailyStats) Run() error {
	for {
		y, m, d := time.Now().Add(time.Hour * 24).Date()
		initTime := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		select {
		case <-ds.stopSig:
			return nil
		case <-time.After(initTime.Sub(time.Now())):
			var stats []dmodels.DailyStat
			for _, act := range ds.actions {
				m, err := act()
				if err != nil {
					actionName := runtime.FuncForPC(reflect.ValueOf(act).Pointer()).Name()
					log.Error("DailyStats: %s: %s", actionName, err.Error())
				}
				for k, v := range m {
					stats = append(stats, dmodels.DailyStat{
						Title:     k,
						Value:     v,
						CreatedAt: initTime,
					})
				}
				err = ds.dao.CreateDailyStats(stats)
				if err != nil {
					actionName := runtime.FuncForPC(reflect.ValueOf(act).Pointer()).Name()
					log.Error("DailyStats: %s: CreateDailyStats: %s", actionName, err.Error())
				}
				log.Info("DailyStats: collection is over, duration: %s", time.Now().Sub(initTime))
			}
		}
	}
}

func (ds *DailyStats) Stop() error {
	ds.stopSig <- struct{}{}
	return nil
}

func (ds *DailyStats) Title() string {
	return "Daily Stats"
}
