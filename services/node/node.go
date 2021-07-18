package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	successfulCode = "successful"

	txInfoByMiniBlockHashEndpoint = "/transactions"
	txInfoByHashEndpoint          = "/transactions/%s"
	miniblockByHashEndpoint       = "/miniblocks/%s"
	blockByNonceAndShardEndpoint  = "/block/%d/by-nonce/%d"
	blockByHashAndShardEndpoint   = "/block/%d/by-hash/%s"
	hyperBlockByNonceEndpoint     = "/hyperblock/by-nonce/%d"
	networkInfoByShardEndpoint    = "/network/status/%d"
	addressEndpoint               = "/address/%s"
	blockExtraDataEndpoint        = "/blocks/%s"
	statsEndpoint                 = "/stats"
	identitiesEndpoint            = "/identities"
	economicsEndpoint             = "/economics"
	networkStatusEndpoint         = "/network/status/%d"
)

type (
	API struct {
		client  *http.Client
		address string
	}

	baseResponse struct {
		Data  json.RawMessage `json:"data"`
		Code  string          `json:"code"`
		Error string          `json:"error"`
	}

	APIi interface {
		GetTxByHash(hash string) (tx Tx, err error)
		GetTxsByMiniBlockHash(miniBlockHash string, offset, limit uint64) (txs []Tx, err error)
		GetMiniBlock(hash string) (miniBlock MiniBlock, err error)
		GetBlock(height uint64, shard uint64) (block Block, err error)
		GetBlockByHash(hash string, shard uint64) (block Block, err error)
		GetHyperBlock(height uint64) (hyperBlock HyperBlock, err error)
		GetMaxHeight(shardIndex uint64) (height uint64, err error)
		GetExtraDataBlock(hash string) (block BlockExtraData, err error)
		GetStats() (stats Stats, err error)
		GetIdentities() (identities []Identity, err error)
		GetEconomics() (economics []Economics, err error)
		GetNetworkStatus(shardID uint64) (status NetworkStatus, err error)
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

func (api *API) GetTxByHash(hash string) (tx Tx, err error) {
	endpoint := fmt.Sprintf(txInfoByHashEndpoint, hash)
	err = api.get(endpoint, nil, &tx, false)
	return tx, err
}

func (api *API) GetTxsByMiniBlockHash(miniBlockHash string, offset, limit uint64) (txs []Tx, err error) {
	params := map[string]string{
		"miniBlockHash": miniBlockHash,
		"from":          fmt.Sprint(offset),
		"size":          fmt.Sprint(limit),
	}
	err = api.get(txInfoByMiniBlockHashEndpoint, params, &txs, false)
	return txs, err
}

func (api *API) GetMiniBlock(hash string) (miniBlock MiniBlock, err error) {
	endpoint := fmt.Sprintf(miniblockByHashEndpoint, hash)
	err = api.get(endpoint, nil, &miniBlock, false)
	return miniBlock, err
}

func (api *API) GetBlock(height uint64, shard uint64) (block Block, err error) {
	endpoint := fmt.Sprintf(blockByNonceAndShardEndpoint, shard, height)
	err = api.get(endpoint, nil, &block, true)
	return block, err
}

func (api *API) GetBlockByHash(hash string, shard uint64) (block Block, err error) {
	endpoint := fmt.Sprintf(blockByHashAndShardEndpoint, shard, hash)
	err = api.get(endpoint, nil, &block, true)
	return block, err
}

func (api *API) GetExtraDataBlock(hash string) (block BlockExtraData, err error) {
	endpoint := fmt.Sprintf(blockExtraDataEndpoint, hash)
	err = api.get(endpoint, nil, &block, false)
	return block, err
}

func (api *API) GetHyperBlock(height uint64) (hyperBlock HyperBlock, err error) {
	endpoint := fmt.Sprintf(hyperBlockByNonceEndpoint, height)
	err = api.get(endpoint, nil, &hyperBlock, true)
	return hyperBlock, err
}

// GetMaxHeight returns current height of specific shard
func (api *API) GetMaxHeight(shardIndex uint64) (height uint64, err error) {
	var chainStatus ChainStatus
	endpoint := fmt.Sprintf(networkInfoByShardEndpoint, shardIndex)
	err = api.get(endpoint, nil, &chainStatus, true)
	return chainStatus.Status.ErdHighestFinalNonce, err
}

func (api *API) GetAddress(address string) (resp Address, err error) {
	endpoint := fmt.Sprintf(addressEndpoint, address)
	err = api.get(endpoint, nil, &resp, true)
	return resp, err
}

func (api *API) GetStats() (stats Stats, err error) {
	err = api.get(statsEndpoint, nil, &stats, false)
	return stats, err
}

func (api *API) GetIdentities() (identities []Identity, err error) {
	err = api.get(identitiesEndpoint, nil, &identities, false)
	return identities, err
}

func (api *API) GetEconomics() (economics []Economics, err error) {
	err = api.get(economicsEndpoint, nil, &economics, false)
	return economics, err
}

func (api *API) GetNetworkStatus(shardID uint64) (status NetworkStatus, err error) {
	endpoint := fmt.Sprintf(networkStatusEndpoint, shardID)
	err = api.get(endpoint, nil, &status, true)
	return status, err
}

func (api *API) get(endpoint string, params map[string]string, result interface{}, wrapped bool) error {
	//<-time.After(time.Millisecond * 200) // todo make latency for tests

	fullURL := fmt.Sprintf("%s%s", api.address, endpoint)
	if len(params) != 0 {
		values := url.Values{}
		for key, value := range params {
			values.Add(key, value)
		}
		fullURL = fmt.Sprintf("%s?%s", fullURL, values.Encode())
	}
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
	finalData := body
	if wrapped {
		var baseResp baseResponse
		err = json.Unmarshal(body, &baseResp)
		if err != nil {
			return fmt.Errorf("json.Unmarshal(baseResponse): %s", err.Error())
		}
		if baseResp.Code != successfulCode {
			return fmt.Errorf("[code: %s] %s", baseResp.Code, baseResp.Error)
		}
		finalData = baseResp.Data
	}
	err = json.Unmarshal(finalData, result)
	if err != nil {
		return fmt.Errorf("json.Unmarshal(result): %s", err.Error())
	}
	return nil
}
