package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetBlock(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["hash"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid hash")
		return
	}
	resp, err := api.svc.GetBlock(address)
	if err != nil {
		log.Error("API GetBlock: svc.GetBlock: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetBlocks(w http.ResponseWriter, r *http.Request) {
	var filter filters.Blocks
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetBlocks: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	resp, err := api.svc.GetBlocks(filter)
	if err != nil {
		log.Error("API GetBlocks: svc.GetBlocks: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetMiniBlock(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["hash"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid hash")
		return
	}
	resp, err := api.svc.GetMiniBlock(address)
	if err != nil {
		log.Error("API GetMiniBlock: svc.GetMiniBlock: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}

