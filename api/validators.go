package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetStakingProvider(w http.ResponseWriter, r *http.Request) {
	address, ok := mux.Vars(r)["address"]
	if !ok || address == "" {
		jsonBadRequest(w, "invalid address")
		return
	}
	provider, err := api.svc.GetStakingProvider(address)
	if err != nil {
		log.Error("API GetStakingProvider: svc.GetStakingProvider: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, provider)
}

func (api *API) GetStakingProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := api.svc.GetStakingProviders()
	if err != nil {
		log.Error("API GetStakingProviders: svc.GetStakingProviders: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, providers)
}

func (api *API) GetNode(w http.ResponseWriter, r *http.Request) {
	key, ok := mux.Vars(r)["key"]
	if !ok || key == "" {
		jsonBadRequest(w, "invalid key")
		return
	}
	provider, err := api.svc.GetNode(key)
	if err != nil {
		log.Error("API GetNode: svc.GetNode: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, provider)
}

func (api *API) GetNodes(w http.ResponseWriter, r *http.Request) {
	var filter filters.Nodes
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetNodes: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(100)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetNodes: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	nodes, err := api.svc.GetNodes(filter)
	if err != nil {
		log.Error("API GetNodes: svc.GetNodes: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, nodes)
}
