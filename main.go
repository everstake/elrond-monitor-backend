package main

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/api"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/services"
	"github.com/everstake/elrond-monitor-backend/services/dailystats"
	"github.com/everstake/elrond-monitor-backend/services/es"
	"github.com/everstake/elrond-monitor-backend/services/modules"
	"github.com/everstake/elrond-monitor-backend/services/parser"
	"github.com/everstake/elrond-monitor-backend/services/scheduler"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	configFilePath = "./config.json"
)

func main() {
	err := os.Setenv("TZ", "UTC")
	if err != nil {
		log.Fatalf("os.Setenv (TZ): %s", err.Error())
	}

	cfg, err := config.GetConfigFromFile(configFilePath)
	if err != nil {
		log.Fatalf("config.GetConfigFromFile: %s", err.Error())
	}

	e, err := es.NewClient(cfg.ElasticSearch.Address)
	if err != nil {
		log.Fatalf("es.NewClient: %s", err.Error())
	}
	//fmt.Println(e.GetBlocks(filters.Blocks{
	//	Pagination: filters.Pagination{
	//		Limit: 10,
	//		Page:  0,
	//	},
	//	Nonce: 5526152,
	//	Shard: []uint64{0},
	//}))

	//txs, err := e.GetTransactions(filters.Transactions{
	//	Pagination: filters.Pagination{
	//		Limit: 10,
	//		Page:  0,
	//	},
	//	MiniBlock: "751a18c1775541958f0af790e4f81ae1623c2532cb2c0f63f398dae1723f3c2d",
	//})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//for _, tx := range txs {
	//	fmt.Println(tx.Hash)
	//}

	fmt.Println(e.GetTransactionsCount(filters.Transactions{
		Pagination: filters.Pagination{
			Limit: 10,
			Page:  0,
		},
		Address: "erd1w4apvcpg7vpkzry0rwyfycpdjquutea4aswwq3pep7khvqaen9lslua6h3",
	}))

	return

	d, err := dao.NewDAO(cfg)
	if err != nil {
		log.Fatalf("dao.NewDAO: %s", err.Error())
	}

	prs := parser.NewParser(cfg, d)

	s, err := services.NewServices(d, cfg, prs)
	if err != nil {
		log.Fatalf("services.NewServices: %s", err.Error())
	}

	ds, err := dailystats.NewDailyStats(cfg, d)
	if err != nil {
		log.Fatalf("dailystats.NewDailyStats: %s", err.Error())
	}

	apiServer := api.NewAPI(cfg, s, d)

	sch := scheduler.NewScheduler()
	sch.AddProcessWithInterval(s.UpdateStats, time.Minute)
	sch.AddProcessWithInterval(s.UpdateValidatorsMap, time.Minute*20)
	sch.AddProcessWithInterval(s.UpdateStakingProviders, time.Hour)
	sch.AddProcessWithInterval(s.UpdateNodes, time.Hour)
	sch.AddProcessWithInterval(s.UpdateValidators, time.Hour)

	g := modules.NewGroup(apiServer, prs, ds, sch)
	g.Run()

	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, os.Interrupt, os.Kill)

	<-gracefulStop
	g.Stop()

	os.Exit(0)
}
