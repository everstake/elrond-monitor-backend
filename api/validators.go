package api

import (
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
