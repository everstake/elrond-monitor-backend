package main

import (
	"log"
	"os"
)

const (
	configFilePath = "./config.json"
)

func main() {
	err := os.Setenv("TZ", "UTC")
	if err != nil {
		log.Fatal("os.Setenv (TZ): %s", err.Error())
	}

	//cfg, err := config.GetConfigFromFile(configFilePath)
	//if err != nil {
	//	log.Fatal("config.GetConfigFromFile: %s", err.Error())
	//}
	//d, err := dao.NewDAO(cfg)
	//if err != nil {
	//	log.Fatal("dao.NewDAO: %s", err.Error())
	//}

	//s, err := services.NewServices(d, cfg)
	//if err != nil {
	//	log.Fatal("services.NewServices: %s", err.Error())
	//}
	//
	//prs := parser.NewParser(cfg, d)
	//
	//apiServer := api.NewAPI(cfg, s, d)

	//sch := scheduler.NewScheduler()

	//g := modules.NewGroup(apiServer, prs)
	//g := modules.NewGroup(apiServer, sch, prs)
	//g.Run()

	//gracefulStop := make(chan os.Signal)
	//signal.Notify(gracefulStop, os.Interrupt, os.Kill)
	//
	//<-gracefulStop
	//g.Stop()
	//
	//os.Exit(0)
}
