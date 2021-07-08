package api

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services"
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
	queryDecoder *schema.Decoder
}

type errResponse struct {
	Error string `json:"error"`
	Msg   string `json:"msg"`
}

func NewAPI(cfg config.Config, svc services.Services, dao dao.DAO) *API {
	sd := schema.NewDecoder()
	sd.IgnoreUnknownKeys(true)
	sd.RegisterConverter(dmodels.Time{}, func(s string) reflect.Value {
		timestamp, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}
		}
		t := dmodels.NewTime(time.Unix(timestamp, 0))
		return reflect.ValueOf(t)
	})
	return &API{
		cfg:          cfg,
		dao:          dao,
		svc:          svc,
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

	api.router.
		PathPrefix("/static").
		Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./resources/static"))))

	wrapper := negroni.New()

	wrapper.Use(cors.New(cors.Options{
		AllowedOrigins:   api.cfg.API.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-User-Env", "Sec-Fetch-Mode"},
	}))

	// public
	HandleActions(api.router, wrapper, "", []*Route{
		{Path: "/", Method: http.MethodGet, Func: api.Index},
		{Path: "/health", Method: http.MethodGet, Func: api.Health},
		{Path: "/api", Method: http.MethodGet, Func: api.GetSwaggerAPI},
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

func jsonError(writer http.ResponseWriter) {
	writer.WriteHeader(500)
	bytes, err := json.Marshal(errResponse{
		Error: "service_error",
	})
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