package es

import (
	"encoding/json"
	"fmt"
	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"io/ioutil"
	"strings"
)

type (
	Client struct {
		cli *elasticsearch.Client
	}

	SearchResponse struct {
		Took int64 `json:"took"`
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []*SearchHit `json:"hits"`
		}
	}
	SearchHit struct {
		Score   float64         `json:"_score"`
		Index   string          `json:"_index"`
		Id      string          `json:"_id"`
		Type    string          `json:"_type"`
		Version int64           `json:"_version,omitempty"`
		Source  json.RawMessage `json:"_source"`
	}
	CountResponse struct {
		Count uint64 `json:"count"`
	}
	obj map[string]interface{}
)

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

func (c *Client) GetBlocks(filter filters.Blocks) (blocks []data.Block, err error) {
	query := obj{
		"sort": obj{
			"timestamp": obj{"order": "desc"},
		},
	}
	if filter.Nonce != 0 {
		addQuery(query, filter.Nonce, "query", "match", "nonce")
	}
	if len(filter.Shard) > 0 {
		addQuery(query, filter.Shard[0], "query", "match", "shard")
	}
	if filter.Limit != 0 {
		query["size"] = filter.Limit
	}
	if filter.Offset() != 0 {
		query["from"] = filter.Offset()
	}
	keys, err := c.search("blocks", query, &blocks)
	if len(keys) != len(blocks) {
		return blocks, fmt.Errorf("wrong number of keys")
	}
	for i, key := range keys {
		blocks[i].Hash = key
	}
	return blocks, err
}

func (c *Client) GetBlocksCount(filter filters.Blocks) (total uint64, err error) {
	query := obj{}
	if filter.Nonce != 0 {
		addQuery(query, filter.Nonce, "query", "match", "nonce")
	}
	if len(filter.Shard) > 0 {
		addQuery(query, filter.Shard[0], "query", "match", "shard")
	}
	total, err = c.count("blocks", query)
	return total, err
}

func (c *Client) GetTransaction(hash string) (tx data.Transaction, err error) {
	err = c.get("transactions", hash, &tx)
	return tx, err
}

func (c *Client) GetTransactions(filter filters.Transactions) (txs []data.Transaction, err error) {
	query := obj{
		"sort": obj{
			"timestamp": obj{"order": "desc"},
		},
	}
	if filter.Address != "" {
		addQuery(query, filter.Address, "query", "multi_match", "query")
		addQuery(query, []string{"sender", "receiver"}, "query", "multi_match", "fields")
	}
	if filter.MiniBlock != "" {
		addQuery(query, filter.MiniBlock, "query", "match", "miniBlockHash")
	}
	if filter.Limit != 0 {
		query["size"] = filter.Limit
	}
	if filter.Offset() != 0 {
		query["from"] = filter.Offset()
	}
	keys, err := c.search("transactions", query, &txs)
	if len(keys) != len(txs) {
		return txs, fmt.Errorf("wrong number of keys")
	}
	for i, key := range keys {
		txs[i].Hash = key
	}
	return txs, err
}

func (c *Client) GetTransactionsCount(filter filters.Transactions) (total uint64, err error) {
	query := obj{}
	if filter.Address != "" {
		addQuery(query, filter.Address, "query", "multi_match", "query")
		addQuery(query, []string{"sender", "receiver"}, "query", "multi_match", "fields")
	}
	if filter.MiniBlock != "" {
		addQuery(query, filter.MiniBlock, "query", "match", "miniBlockHash")
	}
	total, err = c.count("transactions", query)
	return total, err
}

func (c *Client) GetMiniblock(hash string) (miniblock data.Miniblock, err error) {
	err = c.get("miniblocks", hash, &miniblock)
	return miniblock, err
}

func (c *Client) GetMiniblocks(filter filters.MiniBlocks) (txs []data.Miniblock, err error) {
	query := obj{
		"sort": obj{
			"timestamp": obj{"order": "desc"},
		},
	}
	if filter.ParentBlockHash != "" {
		query["bool"] = obj{"must": []obj{
			{"match": obj{"senderBlockHash": filter.ParentBlockHash}},
			{"match": obj{"receiverBlockHash": filter.ParentBlockHash}},
		}}
	}
	keys, err := c.search("miniblocks", query, &txs)
	if len(keys) != len(txs) {
		return txs, fmt.Errorf("wrong number of keys")
	}
	for i, key := range keys {
		txs[i].Hash = key
	}
	return txs, err
}

func (c *Client) GetSCResults(txHash string) (scs []data.ScResult, err error) {
	query := obj{
		"sort": obj{
			"timestamp": obj{"order": "desc"},
		},
		"query": obj{
			"match": obj{
				"originalTxHash": txHash,
			},
		},
	}
	keys, err := c.search("scresults", query, &scs)
	if len(keys) != len(scs) {
		return scs, fmt.Errorf("wrong number of keys")
	}
	for i, key := range keys {
		scs[i].Hash = key
	}
	return scs, err
}

func (c *Client) GetAccount(address string) (acc data.Account, err error) {
	err = c.get("miniblocks", address, &acc)
	return acc, err
}


func (c *Client) GetAccounts(filter filters.Accounts) (accounts []data.Account, err error) {
	query := obj{
		"sort": obj{
			"balanceNum": obj{"order": "desc"},
		},
	}
	if filter.Limit != 0 {
		query["size"] = filter.Limit
	}
	if filter.Offset() != 0 {
		query["from"] = filter.Offset()
	}
	_, err = c.search("accounts", query, &accounts)
	//if len(keys) != len(accounts) {
	//	return accounts, fmt.Errorf("wrong number of keys")
	//}
	//for i, key := range keys {
	//	accounts[i].Hash = key
	//}
	return accounts, err
}


func (c *Client) GetAccountsCount(filter filters.Accounts) (total uint64, err error) {
	total, err = c.count("accounts", nil)
	return total, err
}


func (c *Client) ValidatorsKeys(shard uint64, epoch uint64) (keys data.ValidatorsPublicKeys, err error) {
	err = c.get("validators", fmt.Sprintf("%d_%d", shard, epoch), &keys)
	return keys, err
}

func (c *Client) search(index string, query map[string]interface{}, dst interface{}) (keys []string, err error) {
	resp, err := c.cli.Search(
		c.cli.Search.WithIndex(index),
		c.cli.Search.WithBody(esutil.NewJSONReader(&query)),
	)
	if err != nil {
		return keys, fmt.Errorf("cli.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return keys, fmt.Errorf(resp.String())
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return keys, fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	fmt.Println(string(d))
	var searchResp SearchResponse
	err = json.Unmarshal(d, &searchResp)
	if err != nil {
		return keys, fmt.Errorf("json.Unmarshal(respMap): %s", err.Error())
	}
	items := make([]string, len(searchResp.Hits.Hits))
	for i, hit := range searchResp.Hits.Hits {
		items[i] = string(hit.Source)
		keys = append(keys, hit.Id)
	}
	preparedString := fmt.Sprintf("[%s]", strings.Join(items, ","))
	err = json.Unmarshal([]byte(preparedString), dst)
	if err != nil {
		return keys, fmt.Errorf("json.Unmarshal(preparedString): %s", err.Error())
	}
	return keys, nil
}

func (c *Client) count(index string, query map[string]interface{}) (total uint64, err error) {
	resp, err := c.cli.Count(
		c.cli.Count.WithIndex(index),
		c.cli.Count.WithBody(esutil.NewJSONReader(&query)),
	)
	if err != nil {
		return total, fmt.Errorf("cli.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return total, fmt.Errorf(resp.String())
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return total, fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	var countResp CountResponse
	err = json.Unmarshal(d, &countResp)
	if err != nil {
		return total, fmt.Errorf("json.Unmarshal(CountResponse): %s", err.Error())
	}
	return countResp.Count, nil
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
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	var respMap map[string]json.RawMessage
	err = json.Unmarshal(d, &respMap)
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

func addQuery(query map[string]interface{}, value interface{}, fields ...string) {
	if len(fields) == 0 {
		return
	}
	for i, field := range fields {
		if i == len(fields)-1 {
			query[field] = value
			break
		}
		if _, ok := query[field]; !ok {
			query[field] = make(obj)
		}
		query = query[field].(obj)
	}
}
