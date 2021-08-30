package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"net/http"
	"time"
)

func (api *API) GetStats(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetStats()
	if err != nil {
		log.Error("API GetStats: svc.GetStats: %s", err.Error())
		jsonError(err, w)
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
		if filter.From.IsZero() {
			filter.From = smodels.NewTime(time.Now().Add(-time.Hour * 30))
		}
		filter.Key = key
		resp, err := api.svc.GetDailyStats(filter)
		if err != nil {
			log.Error("API GetStats: svc.GetDailyStats: %s", err.Error())
			jsonError(err, w)
			return
		}
		jsonData(w, resp)
	}
}

func (api *API) GetValidatorsMap(w http.ResponseWriter, r *http.Request) {
	data, err := api.svc.GetValidatorsMap()
	if err != nil {
		log.Error("API GetValidatorsMap: svc.GetValidatorsMap: %s", err.Error())
		jsonError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (api *API) GetValidatorStats(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetValidatorStats()
	if err != nil {
		log.Error("API GetValidatorStats: GetValidatorStats.GetStats: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}
