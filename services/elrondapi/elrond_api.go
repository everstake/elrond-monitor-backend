package elrondapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	statsEndpoint                 = "/stats"
	identitiesEndpoint            = "/identities"
	economicsEndpoint             = "/economics"
	accountDelegationEndpoint     = "/accounts/%s/delegation-legacy"

	host = "https://api.elrond.com"
)

type (
	API struct {
		client  *http.Client
	}

	APIi interface {
		GetStats() (stats Stats, err error)
		GetIdentities() (identities []Identity, err error)
		GetEconomics() (economics []Economics, err error)
		GetAccountDelegation(address string) (account AccountDelegation, err error)
	}
)

func NewAPI() *API {
	return &API{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (api *API) GetStats() (stats Stats, err error) {
	err = api.get(statsEndpoint, nil, &stats)
	return stats, err
}

func (api *API) GetIdentities() (identities []Identity, err error) {
	err = api.get(identitiesEndpoint, nil, &identities)
	return identities, err
}

func (api *API) GetEconomics() (economics []Economics, err error) {
	err = api.get(economicsEndpoint, nil, &economics)
	return economics, err
}

func (api *API) GetAccountDelegation(address string) (account AccountDelegation, err error) {
	endpoint := fmt.Sprintf(accountDelegationEndpoint, address)
	err = api.get(endpoint, nil, &account)
	return account, err
}

func (api *API) get(endpoint string, params map[string]string, result interface{}) error {
	fullURL := fmt.Sprintf("%s%s", host, endpoint)
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
	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("json.Unmarshal(result): %s", err.Error())
	}
	return nil
}
