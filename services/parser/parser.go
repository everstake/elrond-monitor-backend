package parser

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services/es"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/shopspring/decimal"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	repeatDelay     = time.Second * 5
	parserTitle     = "elrond"
	fetcherChBuffer = 5000
	saverChBuffer   = 5000
	msgOKBase64     = "QDZmNmI=" // @ok
	msgOKHex        = "@6f6b"    // @ok
)

type (
	Parser struct {
		cfg       config.Config
		node      node.APIi
		es        *es.Client
		dao       dao.DAO
		fetcherCh chan uint64
		saverCh   chan parsedData
		accounts  map[string]struct{}
		ctx       context.Context
		cancel    context.CancelFunc
		wg        *sync.WaitGroup

		mu          *sync.RWMutex
		delegations map[string]map[string]decimal.Decimal
	}
	parsedData struct {
		height      uint64
		delegations []dmodels.Delegation
		rewards     []dmodels.Reward
		stakeEvents []dmodels.StakeEvent
	}
	ShardIndex uint64
)

func NewParser(cfg config.Config, d dao.DAO) (*Parser, error) {
	esClient, err := es.NewClient(cfg.ElasticSearch.Address)
	if err != nil {
		return nil, fmt.Errorf("es.NewClient: %s", err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Parser{
		cfg:       cfg,
		dao:       d,
		node:      node.NewAPI(cfg.Parser.Node, cfg.Contracts),
		es:        esClient,
		fetcherCh: make(chan uint64, fetcherChBuffer),
		saverCh:   make(chan parsedData, saverChBuffer),
		accounts:  make(map[string]struct{}),
		ctx:       ctx,
		cancel:    cancel,
		wg:        &sync.WaitGroup{},

		mu:          &sync.RWMutex{},
		delegations: make(map[string]map[string]decimal.Decimal),
	}, nil
}

func (p *Parser) Run() error {
	model, err := p.dao.GetParser(parserTitle)
	if err != nil {
		return fmt.Errorf("parser not found")
	}
	err = p.loadStates()
	if err != nil {
		return fmt.Errorf("loadStates: %s", err.Error())
	}
	for i := uint64(0); i < p.cfg.Parser.Fetchers; i++ {
		go p.runFetcher()
	}

	go p.saving()
	for {
		block, err := p.es.GetLatestBlock(node.MetaChainShardIndex)
		if err != nil {
			log.Error("Parser: es.GetLatestBlock: %s", err.Error())
			<-time.After(time.Second)
			continue
		}

		latestBlock := block.Nonce
		if model.Height >= latestBlock {
			<-time.After(time.Second)
			continue
		}
		for ; model.Height < latestBlock; model.Height++ {
			select {
			case <-p.ctx.Done():
				return nil
			case p.fetcherCh <- model.Height + 1:
			}
		}
	}
}

func (p *Parser) Title() string {
	return "Parser"
}

func (p *Parser) Stop() error {
	p.cancel()
	p.wg.Wait()
	return nil
}

func (p *Parser) runFetcher() {
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
		}
		height := <-p.fetcherCh
		for {
			d, err := p.parseHyperBlock(height)
			if err != nil {
				log.Error("Parser: parseHyperBlock(%d): %s", height, err.Error())
				<-time.After(time.Second)
				continue
			}
			p.saverCh <- d
			break
		}

	}
}

func (p *Parser) parseHyperBlock(nonce uint64) (d parsedData, err error) {
	d.height = nonce

	blocks, err := p.es.GetBlocks(filters.Blocks{
		Shard:      []uint64{node.MetaChainShardIndex},
		Nonce:      nonce,
		Pagination: filters.Pagination{Limit: 1},
	})
	if err != nil {
		return d, fmt.Errorf("es.GetBlocks: %s", err.Error())
	}
	if len(blocks) != 1 {
		return d, fmt.Errorf("can`t fetch block")
	}
	hyperBlock := blocks[0]

	hyperBlocks := make([]data.Block, 0)
	hyperBlocks = append(hyperBlocks, hyperBlock)
	for _, bHash := range hyperBlock.NotarizedBlocksHashes {
		block, err := p.es.GetBlock(bHash)
		if err != nil {
			return d, fmt.Errorf("es.GetBlock(%s): %s", bHash, err.Error())
		}
		hyperBlocks = append(hyperBlocks, block)
	}

	for _, block := range hyperBlocks {
		for _, miniBlockHash := range block.MiniBlocksHashes {
			txs, err := p.es.GetTransactions(filters.Transactions{MiniBlock: miniBlockHash})
			if err != nil {
				return d, fmt.Errorf("dao.GetTransactions(miniblock:%s): %s", miniBlockHash, err.Error())
			}

			for _, t := range txs {
				tx, err := p.es.GetTransaction(t.Hash)
				if err != nil {
					return d, fmt.Errorf("es.GetTransaction: %s", err.Error())
				}
				tx.Hash = t.Hash
				if tx.Hash == "6455464a475dc3071b5a6d72965c0157fdd925982c2157f7b46942fb1b683e88" {
					continue
				}

				switch strings.ToLower(tx.Status) {
				case dmodels.TxStatusPending:
					return d, fmt.Errorf("found pending tx: %s", tx.Hash)
				case dmodels.TxStatusSuccess, dmodels.TxStatusFail, dmodels.TxStatusInvalid:
				default:
					return d, fmt.Errorf("unknown tx status: %s", tx.Status)
				}

				if tx.Status != dmodels.TxStatusSuccess {
					continue
				}

				epoch := uint64(block.Epoch)

				switch string(tx.Data) {
				case "withdraw":
					err = d.parseWithdraw(tx, epoch)
					if err != nil {
						return d, fmt.Errorf("[tx_hash: %s] parseWithdraw: %s", tx.Hash, err.Error())
					}
				case "stake":
					err = d.parseStake(tx, epoch)
					if err != nil {
						return d, fmt.Errorf("parseStake: %s", err.Error())
					}
				case "reDelegateRewards": // create delegation + claimRewards
					err = d.parseRewardDelegations(tx, nonce, epoch)
					if err != nil {
						return d, fmt.Errorf("[tx_hash: %s] parseRewardDelegations: %s", tx.Hash, err.Error())
					}
				case "delegate":
					err = d.parseDelegations(tx, epoch)
					if err != nil {
						return d, fmt.Errorf("[tx_hash: %s] parseDelegations: %s", tx.Hash, err.Error())
					}
				case "claimRewards":
					err = d.parseRewardClaims(tx, nonce, epoch)
					if err != nil {
						return d, fmt.Errorf("[tx_hash: %s] parseRewardClaims: %s", tx.Hash, err.Error())
					}
				case "unBondTokens":
					err = d.unBondTokens(tx, epoch)
					if err != nil {
						return d, fmt.Errorf("[tx_hash: %s] unBondTokens: %s", tx.Hash, err.Error())
					}
				default:
					if strings.Contains(string(tx.Data), "unBond") {
						err = d.parseUnbond(tx, epoch)
						if err != nil {
							return d, fmt.Errorf("[tx_hash: %s] parseUnbond: %s", tx.Hash, err.Error())
						}
					}
					if strings.Contains(string(tx.Data), "unDelegate") {
						err = d.parseUndelegations(tx, epoch)
						if err != nil {
							return d, fmt.Errorf("[tx_hash: %s] parseUndelegations: %s", tx.Hash, err.Error())
						}
					}
					if strings.Contains(string(tx.Data), "unStake") {
						err = d.parseUnstake(tx, epoch)
						if err != nil {
							return d, fmt.Errorf("[tx_hash: %s] parseUnstake: %s", tx.Hash, err.Error())
						}
					}
					if strings.Contains(string(tx.Data), "unStakeTokens") {
						err = d.parseUnStakeTokens(tx, epoch)
						if err != nil {
							return d, fmt.Errorf("[tx_hash: %s] unStakeTokens: %s", tx.Hash, err.Error())
						}
					}
				}
			}
		}

	}

	return d, nil
}

func (d *parsedData) parseDelegations(tx es.Tx, epoch uint64) error {
	if !findOK(tx.SCResults) {
		log.Warn("Parser: parseDelegations: findOK: false (tx: %s)", tx.Hash)
		return nil
	}

	amount, err := decimal.NewFromString(tx.Value)
	if err != nil {
		return fmt.Errorf("decimal.NewFromString: %s", err.Error())
	}
	d.delegations = append(d.delegations, dmodels.Delegation{
		Delegator: tx.Sender,
		TxHash:    tx.Hash,
		Validator: tx.Receiver,
		Amount:    node.ValueToEGLD(amount),
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.DelegateStakeEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    node.ValueToEGLD(amount),
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (d *parsedData) parseRewardClaims(tx es.Tx, nonce uint64, epoch uint64) error {
	if len(tx.SCResults) != 2 {
		log.Warn("Parser [tx_hash: %s]: parseRewardClaims: len(tx.ScResults) != 2", tx.Hash)
		return nil
	}
	rewardsIndex := 0
	if string(tx.SCResults[0].Data) == msgOKBase64 || string(tx.SCResults[0].Data) == msgOKHex {
		rewardsIndex = 1
	} else if string(tx.SCResults[1].Data) != msgOKBase64 && string(tx.SCResults[1].Data) != msgOKHex {
		log.Warn("Parser [tx_hash: %s]: parseRewardClaims: can`t find OK msg", tx.Hash)
		return nil
	}
	value, err := decimal.NewFromString(tx.SCResults[rewardsIndex].Value)
	if err != nil {
		return fmt.Errorf("decimal.NewFromString(%s): %s", tx.SCResults[rewardsIndex].Value, err.Error())
	}
	amount := node.ValueToEGLD(value)
	if tooMuchValue(amount) {
		log.Warn("Parser [tx_hash: %s]: parseRewardClaims: too much value", tx.Hash)
		return nil
	}
	d.rewards = append(d.rewards, dmodels.Reward{
		HypeblockID:     nonce,
		TxHash:          tx.Hash,
		ReceiverAddress: tx.Sender,
		Amount:          amount,
		CreatedAt:       time.Unix(int64(tx.Timestamp), 0),
	})
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.ClaimRewardsEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    amount,
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (d *parsedData) parseRewardDelegations(tx es.Tx, nonce uint64, epoch uint64) error {
	if !checkSCResults(tx.SCResults, 2) {
		log.Warn("Parser: parseRewardDelegations: checkSCResults: false (tx: %s)", tx.Hash)
		return nil
	}
	value := tx.SCResults[0].Value
	if len(tx.SCResults[1].Data) == 0 {
		value = tx.SCResults[1].Value
	}
	amount, err := decimal.NewFromString(value)
	if err != nil {
		return fmt.Errorf("decimal.NewFromString(%s): %s", value, err.Error())
	}
	amount = node.ValueToEGLD(amount)
	if tooMuchValue(amount) {
		log.Warn("Parser [tx_hash: %s]: parseRewardDelegations: too much value", tx.Hash)
		return nil
	}
	d.rewards = append(d.rewards, dmodels.Reward{
		HypeblockID:     nonce,
		TxHash:          tx.Hash,
		ReceiverAddress: tx.Sender,
		Amount:          amount,
		CreatedAt:       time.Unix(int64(tx.Timestamp), 0),
	})
	d.delegations = append(d.delegations, dmodels.Delegation{
		Delegator: tx.Sender,
		TxHash:    tx.Hash,
		Validator: tx.Receiver,
		Amount:    amount,
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.ReDelegateRewardsEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    amount,
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (d *parsedData) parseUndelegations(tx es.Tx, epoch uint64) error {
	if !findOK(tx.SCResults) {
		log.Warn("Parser: parseUndelegations: findOK: false (tx: %s)", tx.Hash)
		return nil
	}
	amountData := strings.TrimPrefix(string(tx.Data), "unDelegate@")
	a, err := decimalFromHex(amountData)
	if err != nil {
		log.Warn("Parser [tx_hash: %s]: decimalFromHex: %s", tx.Hash, err.Error())
		return nil
	}
	if tooMuchValue(a) {
		log.Warn("Parser [tx_hash: %s]: parseUndelegations: too much value", tx.Hash)
		return nil
	}
	d.delegations = append(d.delegations, dmodels.Delegation{
		Delegator: tx.Sender,
		TxHash:    tx.Hash,
		Validator: tx.Receiver,
		Amount:    a.Neg(),
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.UnDelegateStakeEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    a.Neg(),
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (d *parsedData) parseStake(tx es.Tx, epoch uint64) error {
	if !findOK(tx.SCResults) {
		log.Warn("Parser: parseStake: findOK: false (tx: %s)", tx.Hash)
		return nil
	}
	amount, err := decimal.NewFromString(tx.Value)
	if err != nil {
		return fmt.Errorf("decimal.NewFromString: %s", err.Error())
	}
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.StakeStakeEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    node.ValueToEGLD(amount),
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (d *parsedData) parseWithdraw(tx es.Tx, epoch uint64) error {
	if !findOK(tx.SCResults) {
		log.Warn("Parser: parseWithdraw: findOK: false (tx: %s)", tx.Hash)
		return nil
	}

	amount := decimal.Zero
	for _, res := range tx.SCResults {
		if tx.Sender == res.Receiver && tx.Receiver == res.Sender && string(res.Data) == "" {
			amount, err := decimal.NewFromString(tx.Value)
			if err != nil {
				return fmt.Errorf("decimal.NewFromString: %s", err.Error())
			}
			amount = node.ValueToEGLD(amount)
			break
		}
	}

	if tooMuchValue(amount) {
		log.Warn("Parser [tx_hash: %s]: parseWithdraw: too much value", tx.Hash)
		return nil
	}
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.WithdrawEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    amount,
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (d *parsedData) parseUnstake(tx es.Tx, epoch uint64) error {
	if !findOK(tx.SCResults) {
		log.Warn("Parser: parseUnstake: findOK: false (tx: %s)", tx.Hash)
		return nil
	}
	amountData := strings.TrimPrefix(string(tx.Data), "unStake@")
	a, err := decimalFromHex(amountData)
	if err != nil {
		log.Warn("Parser [tx_hash: %s]: decimalFromHex: %s", tx.Hash, err.Error())
		return nil
	}
	if tooMuchValue(a) {
		log.Warn("Parser [tx_hash: %s]: parseUnstake: too much value", tx.Hash)
		return nil
	}
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.UnStakeEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    a.Neg(),
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (d *parsedData) parseUnStakeTokens(tx es.Tx, epoch uint64) error {
	if !findOK(tx.SCResults) {
		log.Warn("Parser: parseUnStakeTokens: findOK: false (tx: %s)", tx.Hash)
		return nil
	}
	amountData := strings.TrimPrefix(string(tx.Data), "unStakeTokens@")
	a, err := decimalFromHex(amountData)
	if err != nil {
		log.Warn("Parser [tx_hash: %s]: parseUnStakeTokens: decimalFromHex: %s", tx.Hash, err.Error())
		return nil
	}
	if tooMuchValue(a) {
		log.Warn("Parser [tx_hash: %s]: parseUnStakeTokens: too much value", tx.Hash)
		return nil
	}
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.UnStakeEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    a.Neg(),
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (d *parsedData) parseUnbond(tx es.Tx, epoch uint64) error {
	if len(tx.SCResults) != 2 {
		log.Warn("Parser [tx_hash: %s]: parseUnbond: len SmartContractResults != 2", tx.Hash)
		return nil
	}
	okIndex := 1
	amountIndex := 0
	if string(tx.SCResults[1].Data) == "delegation stake unbond" {
		okIndex = 0
		amountIndex = 1
	} else if string(tx.SCResults[0].Data) != "delegation stake unbond" {
		log.Warn("Parser [tx_hash: %s]: parseUnbond: can`t find `delegation stake unbond`", tx.Hash)
		return nil
	}

	if !strings.Contains(string(tx.SCResults[okIndex].Data), msgOKBase64) && !strings.Contains(string(tx.SCResults[okIndex].Data), msgOKHex) {
		log.Warn("Parser [tx_hash: %s]: parseUnbond: ok not found (%s)`", tx.Hash, tx.SCResults[okIndex].Data)
		return nil
	}

	amount, err := decimal.NewFromString(tx.SCResults[amountIndex].Value)
	if err != nil {
		return fmt.Errorf("decimal.NewFromString: %s", err.Error())
	}

	value := node.ValueToEGLD(amount)
	if tooMuchValue(value) {
		log.Warn("Parser [tx_hash: %s]: parseUnbond: too much value", tx.Hash)
		return nil
	}

	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.UnBondEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    value,
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (d *parsedData) unBondTokens(tx es.Tx, epoch uint64) error {
	if !findOK(tx.SCResults) {
		log.Warn("Parser: unBondTokens: findOK: false (tx: %s)", tx.Hash)
		return nil
	}
	amount := decimal.Zero
	for _, res := range tx.SCResults {
		if tx.Sender == res.Receiver && tx.Receiver == res.Sender && string(res.Data) == "" {
			amount, err := decimal.NewFromString(res.Value)
			if err != nil {
				return fmt.Errorf("decimal.NewFromString: %s", err.Error())
			}
			amount = node.ValueToEGLD(amount)
		}
	}
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    tx.Hash,
		Type:      dmodels.UnBondEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     epoch,
		Amount:    amount,
		CreatedAt: time.Unix(int64(tx.Timestamp), 0),
	})
	return nil
}

func (p *Parser) saving() {
	var model dmodels.Parser
	for {
		var err error
		model, err = p.dao.GetParser(parserTitle)
		if err != nil {
			log.Error("Parser: saving: dao.GetParser: %s", err.Error())
			<-time.After(time.Second * 5)
			continue
		}
		break
	}

	ticker := time.After(time.Second)

	var dataset []parsedData

	for {
		select {
		case <-p.ctx.Done():
			return
		case d := <-p.saverCh:
			dataset = append(dataset, d)
			continue
		case <-ticker:
			sort.Slice(dataset, func(i, j int) bool {
				return dataset[i].height < dataset[j].height
			})
			ticker = time.After(time.Second * 2)
		}

		var count int
		for i, item := range dataset {
			if item.height == model.Height+uint64(i+1) {
				count = i + 1
			} else {
				break
			}
		}

		if count == 0 {
			continue
		}

		if count > int(p.cfg.Parser.Batch) {
			count = int(p.cfg.Parser.Batch)
		}

		var singleData parsedData
		for _, item := range dataset[:count] {
			singleData.delegations = append(singleData.delegations, item.delegations...)
			singleData.rewards = append(singleData.rewards, item.rewards...)
			singleData.stakeEvents = append(singleData.stakeEvents, item.stakeEvents...)
		}
		p.updateStakeStates(singleData.stakeEvents)
		p.wg.Add(1)
		var err error
		for {
			err = p.dao.CreateDelegations(singleData.delegations)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateDelegations: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateRewards(singleData.rewards)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateRewards: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateStakeEvents(singleData.stakeEvents)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateStakeEvents: %s", err.Error())
			<-time.After(repeatDelay)
		}

		for {
			model.Height += uint64(count)
			err = p.dao.UpdateParserHeight(model)
			if err == nil {
				break
			}
			log.Error("Parser: dao.UpdateParserHeight: %s", err.Error())
			<-time.After(repeatDelay)
		}
		dataset = append(dataset[count:])
		p.wg.Done()
	}
}

func (p *Parser) matchMiniblocks(miniblocks []dmodels.MiniBlock) (result []dmodels.MiniBlock) {
	mp := make(map[string]dmodels.MiniBlock)
	for _, mb := range miniblocks {
		b, ok := mp[mb.Hash]
		if !ok {
			mp[mb.Hash] = mb
			continue
		}
		if b.ReceiverBlockHash == "" {
			b.ReceiverBlockHash = mb.ReceiverBlockHash
		}
		if b.SenderBlockHash == "" {
			b.SenderBlockHash = mb.SenderBlockHash
		}
		mp[mb.Hash] = b
	}
	for _, b := range mp {
		result = append(result, b)
	}
	return result
}

func decimalFromHex(hexStr string) (result decimal.Decimal, err error) {
	d, err := hex.DecodeString(hexStr)
	if err != nil {
		return result, fmt.Errorf("hex.DecodeString: %s", err.Error())
	}
	a := (&big.Int{}).SetBytes(d)
	return node.ValueToEGLD(decimal.NewFromBigInt(a, 0)), nil
}

func checkSCResults(results []es.SCResult, expectedLen int) bool {
	if len(results) != expectedLen {
		return false
	}
	switch expectedLen {
	case 1:
		return string(results[0].Data) == msgOKBase64 || string(results[0].Data) == msgOKHex
	case 2:
		okIndex := 0
		if len(results[0].Data) == 0 {
			okIndex = 1
		}
		return string(results[okIndex].Data) == msgOKBase64 || string(results[okIndex].Data) == msgOKHex
	}
	return false
}

func findOK(results []es.SCResult) bool {
	found := false
	for _, res := range results {
		if existOK(string(res.Data)) {
			found = true
			break
		}
	}
	return found
}

func existOK(data string) bool {
	return data == msgOKBase64 || data == msgOKHex
}

func tooMuchValue(d decimal.Decimal) bool {
	return d.GreaterThanOrEqual(decimal.New(1, 18))
}
