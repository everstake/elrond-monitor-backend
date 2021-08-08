package es

import (
	"encoding/json"
	"fmt"
	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/elastic/go-elasticsearch/v7"
	"io/ioutil"
)

type Client struct {
	cli *elasticsearch.Client
}

func NewClient(address string) (*Client, error) {
	cli, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{address},
	})
	if err != nil {
		return nil, fmt.Errorf("elasticsearch.NewClient: %s", err.Error())
	}
	return &Client{cli: cli}, nil
}

func (c *Client) GetBlock(hash string) (block data.Block, err error) {
	err = c.get("blocks", hash, &block)
	return block, err
}

func (c *Client) GetTx(hash string) (tx data.Transaction, err error) {
	err = c.get("transactions", hash, &tx)
	return tx, err
}

func (c *Client) ValidatorsKeys(shard uint64, epoch uint64) (keys data.ValidatorsPublicKeys, err error) {
	err = c.get("validators", fmt.Sprintf("%d_%d", shard, epoch), &keys)
	return keys, err
}

func (c *Client) get(index string, id string, dst interface{}) error {
	resp, err := c.cli.Get(index, id)
	if err != nil {
		return fmt.Errorf("cli.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return fmt.Errorf(resp.String())
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	var respMap map[string]json.RawMessage
	err = json.Unmarshal(data, &respMap)
	if err != nil {
		return fmt.Errorf("json.Unmarshal(respMap): %s", err.Error())
	}
	if _, ok := respMap["_source"]; !ok {
		return fmt.Errorf("not found source | data: %s", resp.String())
	}
	err = json.Unmarshal(respMap["_source"], dst)
	if err != nil {
		return fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}
	return nil
}
