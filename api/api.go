package api

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/api/ws"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services"
	"github.com/everstake/elrond-monitor-backend/services/dailystats"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type API struct {
	dao          dao.DAO
	cfg          config.Config
	svc          services.Services
	router       *mux.Router
	WS           *ws.Hub
	queryDecoder *schema.Decoder
}

type errResponse struct {
	Error string `json:"error"`
	Msg   string `json:"msg,omitempty"`
}

func NewAPI(cfg config.Config, svc services.Services, dao dao.DAO) *API {
	sd := schema.NewDecoder()
	sd.IgnoreUnknownKeys(true)
	sd.RegisterConverter(smodels.Time{}, func(s string) reflect.Value {
		timestamp, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}
		}
		t := smodels.NewTime(time.Unix(timestamp, 0))
		return reflect.ValueOf(t)
	})
	hub := ws.NewHub()
	go hub.Run()
	return &API{
		cfg:          cfg,
		dao:          dao,
		svc:          svc,
		WS:           hub,
		queryDecoder: sd,
	}
}

func (api *API) Title() string {
	return "API"
}

func (api *API) Run() error {
	api.router = mux.NewRouter()
	api.loadRoutes()

	http.Handle("/", api.router)
	log.Info("Listen API server on %d port", api.cfg.API.ListenOnPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", api.cfg.API.ListenOnPort), nil)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) Stop() error {
	return nil
}

func (api *API) loadRoutes() {

	api.router = mux.NewRouter()

	api.router.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(api.WS, w, r)
	})

	api.router.
		PathPrefix("/static").
		Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./resources/static"))))

	wrapper := negroni.New()

	wrapper.Use(cors.New(cors.Options{
		AllowedOrigins:   api.cfg.API.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Sec-Fetch-Mode"},
	}))

	// public
	HandleActions(api.router, wrapper, "", []*Route{
		{Path: "/", Method: http.MethodGet, Func: api.Index},
		{Path: "/health", Method: http.MethodGet, Func: api.Health},
		{Path: "/api", Method: http.MethodGet, Func: api.GetSwaggerAPI},

		{Path: "/transactions", Method: http.MethodGet, Func: api.GetTransactions},
		{Path: "/transaction/{hash}", Method: http.MethodGet, Func: api.GetTransaction},
		{Path: "/blocks", Method: http.MethodGet, Func: api.GetBlocks},
		{Path: "/block/{hash}", Method: http.MethodGet, Func: api.GetBlock},
		{Path: "/block/{shard}/{nonce}", Method: http.MethodGet, Func: api.GetBlockByNonce},
		{Path: "/accounts", Method: http.MethodGet, Func: api.GetAccounts},
		{Path: "/account/{address}", Method: http.MethodGet, Func: api.GetAccount},
		{Path: "/miniblock/{hash}", Method: http.MethodGet, Func: api.GetMiniBlock},
		{Path: "/stats", Method: http.MethodGet, Func: api.GetStats},
		{Path: "/transactions/range", Method: http.MethodGet, Func: api.GetDailyStats(dailystats.TotalTransactionsKey)},
		{Path: "/accounts/range", Method: http.MethodGet, Func: api.GetDailyStats(dailystats.TotalAccountKey)},
		{Path: "/price/range", Method: http.MethodGet, Func: api.GetDailyStats(dailystats.PriceKey)},
		{Path: "/stake/range", Method: http.MethodGet, Func: api.GetDailyStats(dailystats.TotalStakeKey)},
		{Path: "/delegators/range", Method: http.MethodGet, Func: api.GetDailyStats(dailystats.TotalDelegatorsKey)},
		{Path: "/epoch", Method: http.MethodGet, Func: api.GetEpoch},
		{Path: "/validators/map", Method: http.MethodGet, Func: api.GetValidatorsMap},
		{Path: "/stake/events", Method: http.MethodGet, Func: api.GetStakeEvents},
		{Path: "/staking/providers", Method: http.MethodGet, Func: api.GetStakingProviders},
		{Path: "/staking/provider/{address}", Method: http.MethodGet, Func: api.GetStakingProvider},
		{Path: "/nodes", Method: http.MethodGet, Func: api.GetNodes},
		{Path: "/node/{key}", Method: http.MethodGet, Func: api.GetNode},
		{Path: "/validators", Method: http.MethodGet, Func: api.GetValidators},
		{Path: "/validator/{identity}", Method: http.MethodGet, Func: api.GetValidator},
		{Path: "/stats/validators", Method: http.MethodGet, Func: api.GetValidatorStats},
		{Path: "/providers/ranking", Method: http.MethodGet, Func: api.GetRanking},
		{Path: "/operations", Method: http.MethodGet, Func: api.GetOperations},
		{Path: "/esdt/accounts", Method: http.MethodGet, Func: api.GetESDTAccounts},

		// tokens
		{Path: "/token/{identifier}", Method: http.MethodGet, Func: api.GetToken},
		{Path: "/tokens", Method: http.MethodGet, Func: api.GetTokens},
		{Path: "/nft/collection/{identifier}", Method: http.MethodGet, Func: api.GetNFTCollection},
		{Path: "/nft/collections", Method: http.MethodGet, Func: api.GetNFTCollections},
		{Path: "/nft/{identifier}", Method: http.MethodGet, Func: api.GetNFT},
		{Path: "/nfts", Method: http.MethodGet, Func: api.GetNFTs},
	})

}

func jsonData(writer http.ResponseWriter, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("can`t marshal json"))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(bytes)
}

func jsonError(err error, writer http.ResponseWriter) {
	var bytes []byte
	if customErr, ok := err.(smodels.Err); ok {
		writer.WriteHeader(customErr.Code())
		bytes, err = json.Marshal(errResponse{
			Error: customErr.Message(),
		})
	} else {
		writer.WriteHeader(500)
		bytes, err = json.Marshal(errResponse{
			Error: "service_error",
		})
	}
	if err != nil {
		writer.Write([]byte("can`t marshal json"))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(bytes)
}

func jsonBadRequest(writer http.ResponseWriter, msg string) {
	bytes, err := json.Marshal(errResponse{
		Error: "bad_request",
		Msg:   msg,
	})
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("can`t marshal json"))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(400)
	writer.Write(bytes)
}

func (api *API) GetSwaggerAPI(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("./resources/templates/swagger.html")
	if err != nil {
		log.Error("GetSwaggerAPI: ReadFile: ", err)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Error("GetSwaggerAPI: Write: ", err)
		return
	}
}
