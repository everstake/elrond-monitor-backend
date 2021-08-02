package api

import (
	"github.com/everstake/elrond-monitor-backend/log"
	"net/http"
)

func (api *API) GetEpoch(w http.ResponseWriter, r *http.Request) {
	resp, err := api.svc.GetEpoch()
	if err != nil {
		log.Error("API GetEpoch: svc.GetEpoch: %s", err.Error())
		jsonError(w)
		return
	}
	jsonData(w, resp)
}
