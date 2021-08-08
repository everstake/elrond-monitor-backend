package node

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math/big"
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
	vmValuesEndpoint             = "/vm-values/query"
)

var precisionDiv = decimal.New(1, Precision)

type (
	API struct {
		client    *http.Client
		address   string
		contracts config.Contracts
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
		GetUserStake(address string) (us UserStake, err error)
		GetClaimableRewards(address string) (reward decimal.Decimal, err error)
	}
)

func NewAPI(apiAddress string, contracts config.Contracts) *API {
	return &API{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		address:   apiAddress,
		contracts: contracts,
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

func (api *API) GetUserStake(address string) (us UserStake, err error) {
	data, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.Delegation,
		FuncName:  "getUserStakeByType",
		Args:      []string{address},
	})
	if err != nil {
		return us, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(data) < 5 {
		return us, fmt.Errorf("len(data) != 5")
	}
	withdrawOnlyStake, _ := ContractValueToDecimal(data[0])
	waitingStake, _ := ContractValueToDecimal(data[1])
	activeStake, _ := ContractValueToDecimal(data[2])
	unstakedStake, _ := ContractValueToDecimal(data[3])
	deferredPaymentStake, _ := ContractValueToDecimal(data[4])
	return UserStake{
		WithdrawOnlyStake:    withdrawOnlyStake,
		WaitingStake:         waitingStake,
		ActiveStake:          activeStake,
		UnstakedStake:        unstakedStake,
		DeferredPaymentStake: deferredPaymentStake,
	}, err
}

func (api *API) GetClaimableRewards(address string) (reward decimal.Decimal, err error) {
	data, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.Delegation,
		FuncName:  "getClaimableRewards",
		Args:      []string{address},
	})
	if err != nil {
		return reward, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(data) != 1 {
		return reward, fmt.Errorf("len(data) != 1")
	}
	reward, _ = ContractValueToDecimal(data[0])
	return reward, err
}

func (api *API) contractCall(req ContractReq) (data []string, err error) {
	var resp ContractResp
	reqJson, _ := json.Marshal(req)
	err = api.post(vmValuesEndpoint, reqJson, &resp, "data")
	if err != nil {
		return nil, fmt.Errorf("api.post: %s", err.Error())
	}
	if resp.ReturnCode != "ok" {
		return nil, fmt.Errorf("[%s]: %s", resp.ReturnCode, resp.ReturnMessage)
	}
	return resp.ReturnData, nil
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

func (api *API) post(endpoint string, reqBody []byte, respBody interface{}, field string) error {
	fullURL := fmt.Sprintf("%s%s", api.address, endpoint)
	resp, err := api.client.Post(fullURL, "application/json", bytes.NewBuffer(reqBody))
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
	if _, ok := baseResp.Data[field]; !ok {
		return fmt.Errorf("field %s not found", field)
	}
	err = json.Unmarshal(baseResp.Data[field], respBody)
	if err != nil {
		return fmt.Errorf("json.Unmarshal(result): %s", err.Error())
	}
	return nil
}

func ValueToEGLD(value decimal.Decimal) decimal.Decimal {
	return value.Div(precisionDiv)
}

func ContractValueToDecimal(v string) (r decimal.Decimal, err error) {
	base64DecodedData, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return r, fmt.Errorf("base64.DecodeString: %s", err.Error())
	}
	bigInt := (&big.Int{}).SetBytes(base64DecodedData)
	return decimal.NewFromBigInt(bigInt, 0), nil
}

func addressToHex(adr string) (hexAddress string, err error) {
	_, buff, err := bech32.Decode(adr)
	if err != nil {
		return hexAddress, fmt.Errorf("bech32.Decode: %s", err.Error())
	}

	decodedBytes, err := bech32.ConvertBits(buff, 5, 8, false)
	if err != nil {
		return hexAddress, fmt.Errorf("bech32.ConvertBits: %s", err.Error())
	}

	return hex.EncodeToString(decodedBytes), nil
}
