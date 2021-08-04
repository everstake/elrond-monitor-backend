package node

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Precision = 18

	successfulCode = "successful"

	MetaChainShardIndex = 4294967295

	blockByNonceAndShardEndpoint = "/block/%d/by-nonce/%d"
	blockByHashAndShardEndpoint  = "/block/%d/by-hash/%s?withTxs=true"
	hyperBlockByNonceEndpoint    = "/hyperblock/by-nonce/%d"
	transactionEndpoint          = "/transaction/%s?withResults=true"
	addressEndpoint              = "/address/%s"
	networkStatusEndpoint        = "/network/status/%d"
	networkConfigEndpoint        = "/network/config"
	validatorStatisticsEndpoint  = "/validator/statistics"
	heartbeatstatusEndpoint      = "/node/heartbeatstatus"
	networkEconomicsEndpoint     = "/network/economics"
)

var precisionDiv = decimal.New(1, Precision)

type (
	API struct {
		client  *http.Client
		address string
	}

	baseResponse struct {
		Data  map[string]json.RawMessage `json:"data"`
		Code  string                     `json:"code"`
		Error string                     `json:"error"`
	}

	APIi interface {
		GetBlock(height uint64, shard uint64) (block Block, err error)
		GetBlockByHash(hash string, shard uint64) (block Block, err error)
		GetHyperBlock(height uint64) (hyperBlock HyperBlock, err error)
		GetTransaction(hash string) (tx Tx, err error)
		GetAddress(address string) (resp Address, err error)
		GetNetworkStatus(shardID uint64) (status NetworkStatus, err error)
		GetValidatorStatistics() (statistics ValidatorStatistics, err error)
		GetHeartbeatStatus() (status HeartbeatStatus, err error)
		GetNetworkEconomics() (ne NetworkEconomics, err error)
		GetNetworkConfig() (ne NetworkConfig, err error)
	}
)

func NewAPI(apiAddress string) *API {
	return &API{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		address: apiAddress,
	}
}

func (api *API) GetBlock(height uint64, shard uint64) (block Block, err error) {
	endpoint := fmt.Sprintf(blockByNonceAndShardEndpoint, shard, height)
	err = api.get(endpoint, &block, "block")
	return block, err
}

func (api *API) GetBlockByHash(hash string, shard uint64) (block Block, err error) {
	endpoint := fmt.Sprintf(blockByHashAndShardEndpoint, shard, hash)
	err = api.get(endpoint, &block, "block")
	return block, err
}

func (api *API) GetHyperBlock(height uint64) (hyperBlock HyperBlock, err error) {
	endpoint := fmt.Sprintf(hyperBlockByNonceEndpoint, height)
	err = api.get(endpoint, &hyperBlock, "hyperblock")
	return hyperBlock, err
}

func (api *API) GetAddress(address string) (resp Address, err error) {
	endpoint := fmt.Sprintf(addressEndpoint, address)
	err = api.get(endpoint, &resp, "account")
	return resp, err
}

func (api *API) GetNetworkStatus(shardID uint64) (status NetworkStatus, err error) {
	endpoint := fmt.Sprintf(networkStatusEndpoint, shardID)
	err = api.get(endpoint, &status, "status")
	return status, err
}

func (api *API) GetValidatorStatistics() (statistics ValidatorStatistics, err error) {
	err = api.get(validatorStatisticsEndpoint, &statistics, "statistics")
	return statistics, err
}

func (api *API) GetHeartbeatStatus() (status HeartbeatStatus, err error) {
	err = api.get(heartbeatstatusEndpoint, &status, "heartbeats")
	return status, err
}

func (api *API) GetTransaction(hash string) (tx Tx, err error) {
	endpoint := fmt.Sprintf(transactionEndpoint, hash)
	err = api.get(endpoint, &tx, "transaction")
	return tx, err
}

func (api *API) GetNetworkEconomics() (ne NetworkEconomics, err error) {
	err = api.get(networkEconomicsEndpoint, &ne, "metrics")
	return ne, err
}

func (api *API) GetNetworkConfig() (ne NetworkConfig, err error) {
	err = api.get(networkConfigEndpoint, &ne, "config")
	return ne, err
}

func (api *API) get(endpoint string, result interface{}, field string) error {
	fullURL := fmt.Sprintf("%s%s", api.address, endpoint)
	resp, err := api.client.Get(fullURL)
	if err != nil {
		return fmt.Errorf("client.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	var baseResp baseResponse
	err = json.Unmarshal(body, &baseResp)
	if err != nil {
		return fmt.Errorf("json.Unmarshal(baseResponse): %s", err.Error())
	}
	if baseResp.Code != successfulCode {
		return fmt.Errorf("[code: %s] %s", baseResp.Code, baseResp.Error)
	}
	if _, ok := baseResp.Data[field]; !ok {
		return fmt.Errorf("field %s not found", field)
	}
	err = json.Unmarshal(baseResp.Data[field], result)
	if err != nil {
		return fmt.Errorf("json.Unmarshal(result): %s", err.Error())
	}
	return nil
}

func ValueToEGLD(value decimal.Decimal) decimal.Decimal {
	return value.Div(precisionDiv)
}
