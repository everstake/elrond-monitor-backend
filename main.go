package main

import (
	"github.com/everstake/elrond-monitor-backend/api"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/services"
	"github.com/everstake/elrond-monitor-backend/services/dailystats"
	"github.com/everstake/elrond-monitor-backend/services/modules"
	"github.com/everstake/elrond-monitor-backend/services/parser"
	"github.com/everstake/elrond-monitor-backend/services/scheduler"
	"github.com/everstake/elrond-monitor-backend/services/watcher"
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

	d, err := dao.NewDAO(cfg)
	if err != nil {
		log.Fatalf("dao.NewDAO: %s", err.Error())
	}

	prs, err := parser.NewParser(cfg, d)
	if err != nil {
		log.Fatalf("parser.NewParse: %s", err.Error())
	}

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
	sch.AddProcessWithInterval(s.UpdateStats, time.Minute*3)
	sch.AddProcessWithInterval(s.UpdateValidatorsMap, time.Minute*20)
	sch.AddProcessWithInterval(s.UpdateStakingProviders, time.Hour)
	sch.AddProcessWithInterval(s.UpdateNodes, time.Hour)
	sch.AddProcessWithInterval(s.UpdateValidators, time.Hour)
	sch.AddProcessWithInterval(s.MakeRanking, time.Hour)
	sch.AddProcessWithInterval(s.UpdateTokens, time.Hour)

	w := watcher.NewWatcher(d, apiServer.WS)

	g := modules.NewGroup(apiServer, prs, ds, sch, w)
	g.Run()

	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, os.Interrupt, os.Kill)

	<-gracefulStop
	g.Stop()

	os.Exit(0)
}
