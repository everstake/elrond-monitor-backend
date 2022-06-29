package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetTransaction(w http.ResponseWriter, r *http.Request) {
	hash, ok := mux.Vars(r)["hash"]
	if !ok || hash == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	resp, err := api.svc.GetTransaction(hash)
	if err != nil {
		log.Error("API GetTransaction: svc.GetTransaction: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var filter filters.Transactions
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetTransactions: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(100)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetTransactions: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetTransactions(filter)
	if err != nil {
		log.Error("API GetTransactions: svc.GetTransactions: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetOperations(w http.ResponseWriter, r *http.Request) {
	var filter filters.Operations
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetOperations: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(100)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetOperations: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetOperations(filter)
	if err != nil {
		log.Error("API GetOperations: svc.GetOperations: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}
