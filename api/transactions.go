package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetTransaction(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	resp, err := api.svc.GetTransaction(address)
	if err != nil {
		log.Error("API GetTransaction: svc.GetTransaction: %s", err.Error())
		jsonError(w)
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
	resp, err := api.svc.GetTransactions(filter)
	if err != nil {
		log.Error("API GetTransactions: svc.GetTransactions: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}
