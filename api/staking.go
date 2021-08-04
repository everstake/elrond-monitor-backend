package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"net/http"
)

func (api *API) GetStakeEvents(w http.ResponseWriter, r *http.Request) {
	var filter filters.StakeEvents
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetStakeEvents: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(100)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetStakeEvents: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetStakeEvents(filter)
	if err != nil {
		log.Error("API GetStakeEvents: svc.GetStakeEvents: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}
