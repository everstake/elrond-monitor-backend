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
	var filter filters.StakingProviders
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetStakingProviders: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(5000)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetStakingProviders: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	providers, err := api.svc.GetStakingProviders(filter)
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

func (api *API) GetValidator(w http.ResponseWriter, r *http.Request) {
	identity, ok := mux.Vars(r)["identity"]
	if !ok || identity == "" {
		jsonBadRequest(w, "invalid identity")
		return
	}
	provider, err := api.svc.GetValidator(identity)
	if err != nil {
		log.Error("API GetValidator: svc.GetValidator: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, provider)
}

func (api *API) GetValidators(w http.ResponseWriter, r *http.Request) {
	var filter filters.Validators
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetValidators: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(5000)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetValidators: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	validators, err := api.svc.GetValidators(filter)
	if err != nil {
		log.Error("API GetValidators: svc.GetValidators: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, validators)
}

func (api *API) GetRanking(w http.ResponseWriter, r *http.Request) {
	ranking, err := api.svc.GetRanking()
	if err != nil {
		log.Error("API GetRanking: svc.GetRanking: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, ranking)
}
