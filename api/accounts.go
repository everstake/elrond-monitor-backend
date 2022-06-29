package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/gorilla/mux"
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
	filter.SetMaxLimit(100)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetAccounts: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetAccounts(filter)
	if err != nil {
		log.Error("API GetAccounts: svc.GetAccounts: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetAccount(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok || address == "" || len(address) != 62 {
		jsonBadRequest(w, "invalid address")
		return
	}
	resp, err := api.svc.GetAccount(address)
	if err != nil {
		log.Error("API GetAccount: svc.GetAccount: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetESDTAccounts(w http.ResponseWriter, r *http.Request) {
	var filter filters.ESDT
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetESDTAccounts: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(100)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetESDTAccounts: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetESDTAccounts(filter)
	if err != nil {
		log.Error("API GetESDTAccounts: svc.GetESDTAccounts: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}
