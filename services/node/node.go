package node

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
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
	esdtsEndpoint                = "/network/esdts"
	fungibleESDTEndpoint         = "/network/esdt/fungible-tokens"
	esdtSupplyEndpoint           = "/network/esdt/supply/%s"
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
		GetHeartbeatStatus() (status []HeartbeatStatus, err error)
		GetNetworkEconomics() (ne NetworkEconomics, err error)
		GetNetworkConfig() (ne NetworkConfig, err error)
		GetUserStake(address string) (us UserStake, err error)
		GetClaimableRewards(address string) (reward decimal.Decimal, err error)
		GetProviderAddresses() (addresses []string, err error)
		GetProviderConfig(provider string) (config ProviderConfig, err error)
		GetProviderMeta(provider string) (config ProviderMeta, err error)
		GetProviderNumUsers(provider string) (count uint64, err error)
		GetCumulatedRewards(provider string) (amount decimal.Decimal, err error)
		GetQueue() (items []QueueItem, err error)
		GetOwner(bls string) (owner string, err error)
		GetTotalStakedTopUpStakedBlsKeys(address string) (stake StakeTopup, err error)
		GetESDTProperties(identifier string) (prop ESDTProperties, err error)
		GetESDTAllAddressesAndRoles(identifier string) (addresses []AddressAndRoles, err error)
		GetESDTs() (tokens []string, err error)
		GetFungibleESDTs() (tokens []string, err error)
		GetESDTSupply(ident string) (supply decimal.Decimal, err error)
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

func (api *API) GetHeartbeatStatus() (status []HeartbeatStatus, err error) {
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

func (api *API) GetESDTs() (tokens []string, err error) {
	err = api.get(esdtsEndpoint, &tokens, "tokens")
	return tokens, err
}

func (api *API) GetFungibleESDTs() (tokens []string, err error) {
	err = api.get(fungibleESDTEndpoint, &tokens, "tokens")
	return tokens, err
}

func (api *API) GetESDTSupply(ident string) (supply decimal.Decimal, err error) {
	err = api.get(fmt.Sprintf(esdtSupplyEndpoint, ident), &supply, "supply")
	return supply, err
}

func (api *API) GetUserStake(address string) (us UserStake, err error) {
	hexAddress, err := addressToHex(address)
	if err != nil {
		return us, fmt.Errorf("addressToHex: %s", err.Error())
	}
	data, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.Delegation,
		FuncName:  "getUserStakeByType",
		Args:      []string{hexAddress},
	})
	if err != nil {
		return us, fmt.Errorf("user: %s, contractCall: %s", address, err.Error())
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
	}, nil
}

func (api *API) GetClaimableRewards(address string) (reward decimal.Decimal, err error) {
	hexAddress, err := addressToHex(address)
	if err != nil {
		return reward, fmt.Errorf("addressToHex: %s", err.Error())
	}
	data, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.Delegation,
		FuncName:  "getClaimableRewards",
		Args:      []string{hexAddress},
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

func (api *API) GetProviderAddresses() (addresses []string, err error) {
	rows, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.DelegationManager,
		FuncName:  "getAllContractAddresses",
	})
	if err != nil {
		return nil, fmt.Errorf("contractCall: %s", err.Error())
	}
	addresses = make([]string, len(rows))
	for i, row := range rows {
		addresses[i], err = base64ToAddress(row)
		if err != nil {
			return nil, fmt.Errorf("base64ToAddress: %s", err.Error())
		}
	}
	return addresses, nil
}

func (api *API) GetESDTProperties(identifier string) (prop ESDTProperties, err error) {
	hexIdentifier := hex.EncodeToString([]byte(identifier))
	rows, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.ESDTContract,
		FuncName:  "getTokenProperties",
		Args:      []string{hexIdentifier},
	})
	if err != nil {
		return prop, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(rows) != 18 {
		return prop, fmt.Errorf("wrong count of rows, got %d", len(rows))
	}
	decimalsStr := strings.TrimPrefix(mustB64Decode(rows[5]), "NumDecimals-")
	decimals, err := strconv.ParseUint(decimalsStr, 10, 64)
	if err != nil {
		return prop, fmt.Errorf("wrong decimals %s", decimalsStr)
	}
	getBool := func(text string, prefix string) bool {
		if strings.TrimPrefix(mustB64Decode(text), prefix) == "true" {
			return true
		}
		return false
	}
	owner, err := base64ToAddress(rows[2])
	if err != nil {
		return prop, fmt.Errorf("base64ToAddress: %s", err.Error())
	}
	return ESDTProperties{
		Name:                     mustB64Decode(rows[0]),
		Type:                     mustB64Decode(rows[1]),
		Owner:                    owner,
		Decimals:                 uint(decimals),
		IsPaused:                 getBool(rows[6], "IsPaused-"),
		CanUpgrade:               getBool(rows[7], "CanUpgrade-"),
		CanMint:                  getBool(rows[8], "CanMint-"),
		CanBurn:                  getBool(rows[9], "CanBurn-"),
		CanChangeOwner:           getBool(rows[10], "CanChangeOwner-"),
		CanPause:                 getBool(rows[11], "CanPause-"),
		CanFreeze:                getBool(rows[12], "CanFreeze-"),
		CanWipe:                  getBool(rows[13], "CanWipe-"),
		CanAddSpecialRoles:       getBool(rows[14], "CanAddSpecialRoles-"),
		CanTransferNFTCreateRole: getBool(rows[15], "canTransferNFTCreateRole-"),
		NFTCreateStopped:         getBool(rows[16], "NFTCreateStopped-"),
		Wiped:                    getBool(rows[17], "Wiped-"),
	}, nil
}

func (api *API) GetESDTAllAddressesAndRoles(identifier string) (addresses []AddressAndRoles, err error) {
	hexIdentifier := hex.EncodeToString([]byte(identifier))
	rows, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.ESDTContract,
		FuncName:  "getAllAddressesAndRoles",
		Args:      []string{hexIdentifier},
	})
	if err != nil {
		return addresses, fmt.Errorf("contractCall: %s", err.Error())
	}
	for _, row := range rows {
		if len(row) == 44 {
			address, err := base64ToAddress(row)
			if err != nil {
				return addresses, fmt.Errorf("can`t convert %s to address: %s", row, err.Error())
			}
			addresses = append(addresses, AddressAndRoles{
				Address: address,
			})
			continue
		}
		if len(addresses) == 0 {
			return addresses, errors.New("incorrect response, address must be first")
		}
		role := mustB64Decode(row)
		addresses[len(addresses)-1].Roles = append(addresses[len(addresses)-1].Roles, role)
	}
	return addresses, nil
}

func (api *API) GetProviderConfig(provider string) (config ProviderConfig, err error) {
	rows, err := api.contractCall(ContractReq{
		SCAddress: provider,
		FuncName:  "getContractConfig",
	})
	if err != nil {
		return config, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(rows) < 3 {
		return config, fmt.Errorf("len(rows) < 3")
	}
	owner, err := base64ToAddress(rows[0])
	if err != nil {
		return config, fmt.Errorf("base64ToAddress: %s", err.Error())
	}
	serviceFee, err := ContractValueToDecimal(rows[1])
	if err != nil {
		return config, fmt.Errorf("ContractValueToDecimal: %s", err.Error())
	}
	delegationCap, err := ContractValueToDecimal(rows[2])
	if err != nil {
		return config, fmt.Errorf("ContractValueToDecimal: %s", err.Error())
	}
	return ProviderConfig{
		Owner:         owner,
		ServiceFee:    serviceFee,
		DelegationCap: delegationCap,
	}, nil
}

func (api *API) GetProviderMeta(provider string) (config ProviderMeta, err error) {
	rows, err := api.contractCall(ContractReq{
		SCAddress: provider,
		FuncName:  "getMetaData",
	})
	if err != nil {
		if err.Error() == "[user error]: delegation meta data is not set" {
			return config, nil
		}
		return config, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(rows) != 3 {
		return config, fmt.Errorf("len(rows) != 3")
	}
	return ProviderMeta{
		Name:      mustB64Decode(rows[0]),
		Website:   mustB64Decode(rows[1]),
		Iidentity: mustB64Decode(rows[2]),
	}, nil
}

func (api *API) GetProviderNumUsers(provider string) (count uint64, err error) {
	rows, err := api.contractCall(ContractReq{
		SCAddress: provider,
		FuncName:  "getNumUsers",
	})
	if err != nil {
		return count, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(rows) != 1 {
		return count, fmt.Errorf("len(rows) != 1")
	}
	v, err := ContractValueToDecimal(rows[0])
	if err != nil {
		return count, fmt.Errorf("ContractValueToDecimal: %s", err.Error())
	}
	return v.BigInt().Uint64(), nil
}

func (api *API) GetCumulatedRewards(provider string) (amount decimal.Decimal, err error) {
	rows, err := api.contractCall(ContractReq{
		SCAddress: provider,
		FuncName:  "getTotalCumulatedRewards",
		Caller:    "erd1qqqqqqqqqqqqqqqpqqqqqqqqlllllllllllllllllllllllllllsr9gav8",
	})
	if err != nil {
		return amount, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(rows) != 1 {
		return amount, fmt.Errorf("len(rows) != 1")
	}
	amount, err = ContractValueToDecimal(rows[0])
	if err != nil {
		return amount, fmt.Errorf("ContractValueToDecimal: %s", err.Error())
	}
	return amount, nil
}

func (api *API) GetUserActiveStake(account, provider string) (amount decimal.Decimal, err error) {
	hexAddress, err := addressToHex(account)
	if err != nil {
		return amount, fmt.Errorf("addressToHex: %s", err.Error())
	}
	rows, err := api.contractCall(ContractReq{
		SCAddress: provider,
		FuncName:  "getUserActiveStake",
		Args:      []string{hexAddress},
	})
	if err != nil {
		return amount, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(rows) != 1 {
		return amount, fmt.Errorf("len(rows) != 1")
	}
	amount, err = ContractValueToDecimal(rows[0])
	if err != nil {
		return amount, fmt.Errorf("ContractValueToDecimal: %s", err.Error())
	}
	return amount, nil
}

func (api *API) GetQueue() (items []QueueItem, err error) {
	rows, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.Staking,
		FuncName:  "getQueueRegisterNonceAndRewardAddress",
		Caller:    api.contracts.Auction,
	})
	if err != nil {
		if err.Error() == "[user error]: no one in waitingList" {
			return nil, nil
		}
		return nil, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(rows) == 0 {
		return nil, nil
	}
	if len(rows)%3 != 0 {
		return nil, fmt.Errorf("wrong len of returned data")
	}
	for i := 2; i < len(rows); i += 3 {
		bls := hex.EncodeToString([]byte(mustB64Decode(rows[i-2])))
		rewardsAdr, err := base64ToAddress(rows[i-1])
		if err != nil {
			return nil, fmt.Errorf("ContractValueToDecimal: %s", err.Error())
		}
		nonce, err := ContractValueToDecimal(rows[i-2])
		if err != nil {
			return nil, fmt.Errorf("ContractValueToDecimal: %s", err.Error())
		}
		items = append(items, QueueItem{
			BLS:      bls,
			Provider: rewardsAdr,
			Nonce:    nonce.BigInt().Uint64(),
			Position: int64((i + 1) / 3),
		})
	}

	return items, nil
}

func (api *API) GetOwner(bls string) (owner string, err error) {
	rows, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.Staking,
		FuncName:  "getOwner",
		Caller:    api.contracts.Auction,
		Args:      []string{bls},
	})
	if err != nil {
		return owner, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(rows) != 1 {
		return owner, nil
	}
	owner, err = base64ToAddress(rows[0])
	if err != nil {
		return owner, fmt.Errorf("base64ToAddress: %s", err.Error())
	}
	return owner, nil
}

func (api *API) GetTotalStakedTopUpStakedBlsKeys(address string) (stake StakeTopup, err error) {
	hexAddress, err := addressToHex(address)
	if err != nil {
		return stake, fmt.Errorf("addressToHex: %s", err.Error())
	}
	rows, err := api.contractCall(ContractReq{
		SCAddress: api.contracts.Auction,
		FuncName:  "getTotalStakedTopUpStakedBlsKeys",
		Caller:    api.contracts.Auction,
		Args:      []string{hexAddress},
	})
	if err != nil {
		return stake, fmt.Errorf("contractCall: %s", err.Error())
	}
	if len(rows) == 0 {
		return stake, nil
	}
	if len(rows) < 4 {
		return stake, fmt.Errorf("worong len of rows")
	}
	stake.TopUp, err = ContractValueToDecimal(rows[0])
	if err != nil {
		return stake, fmt.Errorf("ContractValueToDecimal(topUp): %s", err.Error())
	}
	stake.Stake, err = ContractValueToDecimal(rows[1])
	if err != nil {
		return stake, fmt.Errorf("ContractValueToDecimal(staked): %s", err.Error())
	}
	stake.Stake = stake.Stake.Sub(stake.TopUp)
	numNodes, err := ContractValueToDecimal(rows[2])
	if err != nil {
		return stake, fmt.Errorf("ContractValueToDecimal(numNodes): %s", err.Error())
	}
	stake.NumNodes = numNodes.BigInt().Uint64()
	if stake.Stake.Equal(decimal.Zero) && stake.NumNodes == 0 {
		stake.TopUp = decimal.Zero
	}
	for _, b := range rows[3:] {
		bts, _ := base64.StdEncoding.DecodeString(b)
		stake.Blses = append(stake.Blses, hex.EncodeToString(bts))
	}
	stake.Address = address
	stake.Locked = stake.Stake.Add(stake.TopUp)
	if stake.NumNodes > 0 {
		stake.Stake = stake.Stake.Div(decimal.New(int64(stake.NumNodes), 0))
		stake.TopUp = stake.TopUp.Div(decimal.New(int64(stake.NumNodes), 0))
		stake.Locked = stake.Locked.Div(decimal.New(int64(stake.NumNodes), 0))
	}
	return stake, nil
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

func base64ToAddress(b64 string) (address string, err error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return address, fmt.Errorf("base64.DecodeString: %s", err.Error())
	}
	conv, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return address, fmt.Errorf("bech32.ConvertBits: %s", err.Error())
	}
	converted, err := bech32.Encode("erd", conv)
	if err != nil {
		return address, fmt.Errorf("bech32.Encode: %s", err.Error())
	}
	return converted, nil
}

func mustB64Decode(b64 string) string {
	r, _ := base64.StdEncoding.DecodeString(b64)
	return string(r)
}
