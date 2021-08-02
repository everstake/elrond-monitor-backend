package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"net/http"
)

func (api *API) GetStats(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetStats()
	if err != nil {
		log.Error("API GetStats: svc.GetStats: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetDailyStats(key string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var filter filters.DailyStats
		err := api.queryDecoder.Decode(&filter, r.URL.Query())
		if err != nil {
			log.Debug("API GetDailyStats: Decode: %s", err.Error())
			jsonBadRequest(w, "bad params")
			return
		}
		if filter.Limit == 0 {
			filter.Limit = 30
		}
		filter.Key = key
		resp, err := api.svc.GetDailyStats(filter)
		if err != nil {
			log.Error("API GetStats: svc.GetDailyStats: %s", err.Error())
			jsonError(w)
			return
		}
		jsonData(w, resp)
	}
}
