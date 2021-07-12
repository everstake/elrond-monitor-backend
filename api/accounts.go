package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"net/http"
)

func (api *API) GetAccounts(w http.ResponseWriter, r *http.Request) {
	var filter filters.Accounts
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetAccounts: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	resp, err := api.svc.GetAccounts(filter)
	if err != nil {
		log.Error("API GetAccounts: svc.GetAccounts: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}
