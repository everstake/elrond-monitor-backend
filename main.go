package main

import (
	"github.com/everstake/elrond-monitor-backend/api"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/services"
	"github.com/everstake/elrond-monitor-backend/services/dailystats"
	"github.com/everstake/elrond-monitor-backend/services/modules"
	"github.com/everstake/elrond-monitor-backend/services/parser"
	"log"
	"os"
	"os/signal"
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

	s, err := services.NewServices(d, cfg)
	if err != nil {
		log.Fatalf("services.NewServices: %s", err.Error())
	}

	ds := dailystats.NewDailyStats(cfg, d)

	prs := parser.NewParser(cfg, d)

	apiServer := api.NewAPI(cfg, s, d)

	g := modules.NewGroup(apiServer, prs, ds)
	g.Run()

	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, os.Interrupt, os.Kill)

	<-gracefulStop
	g.Stop()

	os.Exit(0)
}
