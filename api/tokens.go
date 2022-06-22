package api

import (
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) GetToken(w http.ResponseWriter, r *http.Request) {
	identifier, ok := mux.Vars(r)["identifier"]
	if !ok || identifier == "" {
		jsonBadRequest(w, "invalid identifier")
		return
	}
	resp, err := api.svc.GetToken(identifier)
	if err != nil {
		log.Error("API GetToken: svc.GetToken: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetTokens(w http.ResponseWriter, r *http.Request) {
	var filter filters.Tokens
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetTokens: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(100)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetTokens: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetTokens(filter)
	if err != nil {
		log.Error("API GetTokens: svc.GetTokens: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetNFTCollection(w http.ResponseWriter, r *http.Request) {
	identifier, ok := mux.Vars(r)["identifier"]
	if !ok || identifier == "" {
		jsonBadRequest(w, "invalid identifier")
		return
	}
	resp, err := api.svc.GetNFTCollection(identifier)
	if err != nil {
		log.Error("API GetNFTCollection: svc.GetNFTCollection: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetNFTCollections(w http.ResponseWriter, r *http.Request) {
	var filter filters.NFTCollections
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetNFTCollections: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(100)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetNFTCollections: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetNFTCollections(filter)
	if err != nil {
		log.Error("API GetNFTCollections: svc.GetNFTCollections: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetNFT(w http.ResponseWriter, r *http.Request) {
	identifier, ok := mux.Vars(r)["identifier"]
	if !ok || identifier == "" {
		jsonBadRequest(w, "invalid identifier")
		return
	}
	resp, err := api.svc.GetNFT(identifier)
	if err != nil {
		log.Error("API GetNFT: svc.GetNFT: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}

func (api *API) GetNFTs(w http.ResponseWriter, r *http.Request) {
	var filter filters.NFTTokens
	err := api.queryDecoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Debug("API GetNFTs: Decode: %s", err.Error())
		jsonBadRequest(w, "bad params")
		return
	}
	filter.SetMaxLimit(100)
	err = filter.Validate()
	if err != nil {
		log.Debug("API GetNFTs: filter.Validate: %s", err.Error())
		jsonBadRequest(w, err.Error())
		return
	}
	resp, err := api.svc.GetNFTs(filter)
	if err != nil {
		log.Error("API GetNFTs: svc.GetNFTs: %s", err.Error())
		jsonError(err, w)
		return
	}
	jsonData(w, resp)
}
