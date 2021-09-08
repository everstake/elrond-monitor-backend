package parser

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/log"
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
		dao       dao.DAO
		fetcherCh chan uint64
		saverCh   chan data
		accounts  map[string]struct{}
		ctx       context.Context
		cancel    context.CancelFunc
		wg        *sync.WaitGroup

		mu          *sync.RWMutex
		delegations map[string]map[string]decimal.Decimal
	}
	data struct {
		height      uint64
		delegations []dmodels.Delegation
		rewards     []dmodels.Reward
		stakeEvents []dmodels.StakeEvent
	}
	ShardIndex uint64
)

func NewParser(cfg config.Config, d dao.DAO) *Parser {
	ctx, cancel := context.WithCancel(context.Background())
	return &Parser{
		cfg:       cfg,
		dao:       d,
		node:      node.NewAPI(cfg.Parser.Node, cfg.Contracts),
		fetcherCh: make(chan uint64, fetcherChBuffer),
		saverCh:   make(chan data, saverChBuffer),
		accounts:  make(map[string]struct{}),
		ctx:       ctx,
		cancel:    cancel,
		wg:        &sync.WaitGroup{},

		mu:          &sync.RWMutex{},
		delegations: make(map[string]map[string]decimal.Decimal),
	}
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
		networkStatus, err := p.node.GetNetworkStatus(node.MetaChainShardIndex)
		if err != nil {
			log.Error("Parser: node.GetMaxHeight: %s", err.Error())
			<-time.After(time.Second)
			continue
		}
		latestBlock := networkStatus.ErdNonce
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

func (p *Parser) parseHyperBlock(nonce uint64) (d data, err error) {
	d.height = nonce

	hyperBlock, err := p.node.GetHyperBlock(nonce)
	if err != nil {
		return d, fmt.Errorf("node.GetHyperBlock: %s", err.Error())
	}

	hyperBlocks := make([]node.Block, 0)
	metaChainBlock, err := p.node.GetBlockByHash(hyperBlock.Hash, node.MetaChainShardIndex)
	if err != nil {
		return d, fmt.Errorf("api.GetBlockByHash(%s): %s", hyperBlock.Hash, err.Error())
	}
	hyperBlocks = append(hyperBlocks, metaChainBlock)
	for _, ShardBlockInfo := range hyperBlock.Shardblocks {
		block, err := p.node.GetBlockByHash(ShardBlockInfo.Hash, ShardBlockInfo.Shard)
		if err != nil {
			return d, fmt.Errorf("api.GetBlockByHash(%s): %s", ShardBlockInfo.Hash, err.Error())
		}
		hyperBlocks = append(hyperBlocks, block)
	}

	for _, block := range hyperBlocks {
		t := time.Unix(block.Timestamp, 0)

		for _, miniBlock := range block.Miniblocks {
			for _, mbTx := range miniBlock.Transactions {
				tx, err := p.node.GetTransaction(mbTx.Hash)
				if err != nil {
					return d, fmt.Errorf("node.GetTransaction(%s): %s", mbTx.Hash, err.Error())
				}

				switch strings.ToLower(tx.Status) {
				case dmodels.TxStatusPending:
					return d, fmt.Errorf("found pending tx: %s", mbTx.Hash)
				case dmodels.TxStatusSuccess, dmodels.TxStatusFail, dmodels.TxStatusInvalid:
				default:
					return d, fmt.Errorf("unknown tx status: %s", tx.Status)
				}

				decodedBytes, err := base64.StdEncoding.DecodeString(tx.Data)
				if err != nil {
					return d, fmt.Errorf("base64.DecodeString: %s", err.Error())
				}

				if tx.Status != dmodels.TxStatusSuccess {
					continue
				}

				txType := string(decodedBytes)
				switch txType {
				case "withdraw":
					err = d.parseWithdraw(tx, mbTx.Hash, t)
					if err != nil {
						return d, fmt.Errorf("[tx_hash: %s] parseWithdraw: %s", mbTx.Hash, err.Error())
					}
				case "stake":
					err = d.parseStake(tx, mbTx.Hash, t)
					if err != nil {
						return d, fmt.Errorf("parseStake: %s", err.Error())
					}
				case "reDelegateRewards": // create delegation + claimRewards
					err = d.parseRewardDelegations(tx, mbTx.Hash, nonce, t)
					if err != nil {
						return d, fmt.Errorf("[tx_hash: %s] parseRewardDelegations: %s", mbTx.Hash, err.Error())
					}
				case "delegate":
					err = d.parseDelegations(tx, mbTx.Hash, t)
					if err != nil {
						return d, fmt.Errorf("[tx_hash: %s] parseDelegations: %s", mbTx.Hash, err.Error())
					}
				case "claimRewards":
					err = d.parseRewardClaims(tx, mbTx.Hash, nonce, t)
					if err != nil {
						return d, fmt.Errorf("[tx_hash: %s] parseRewardClaims: %s", mbTx.Hash, err.Error())
					}
				default:
					if strings.Contains(txType, "unBond") {
						err = d.parseUnbond(tx, mbTx.Hash, t)
						if err != nil {
							return d, fmt.Errorf("[tx_hash: %s] parseUnbond: %s", mbTx.Hash, err.Error())
						}
					}
					if strings.Contains(txType, "unDelegate") {
						err = d.parseUndelegations(tx, mbTx.Hash, txType, t)
						if err != nil {
							return d, fmt.Errorf("[tx_hash: %s] parseUndelegations: %s", mbTx.Hash, err.Error())
						}
					}
					if strings.Contains(txType, "unStake") {
						err = d.parseUnstake(tx, mbTx.Hash, txType, t)
						if err != nil {
							return d, fmt.Errorf("[tx_hash: %s] parseUnstake: %s", mbTx.Hash, err.Error())
						}
					}
				}
			}
		}

	}

	return d, nil
}

func (d *data) parseDelegations(tx node.Tx, txHash string, t time.Time) error {
	if !checkSCResults(tx.SmartContractResults, 2) {
		log.Warn("Parser: parseDelegations: checkSCResults: false (tx: %s)", txHash)
		return nil
	}
	amount, err := decimal.NewFromString(tx.Value)
	if err != nil {
		return fmt.Errorf("decimal.NewFromString: %s", err.Error())
	}
	d.delegations = append(d.delegations, dmodels.Delegation{
		Delegator: tx.Sender,
		TxHash:    txHash,
		Validator: tx.Receiver,
		Amount:    node.ValueToEGLD(amount),
		CreatedAt: t,
	})
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    txHash,
		Type:      dmodels.DelegateStakeEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     tx.Epoch,
		Amount:    node.ValueToEGLD(amount),
		CreatedAt: t,
	})
	return nil
}

func (d *data) parseRewardClaims(tx node.Tx, txHash string, nonce uint64, t time.Time) error {
	if len(tx.SmartContractResults) != 2 {
		log.Warn("Parser [tx_has: %s]: parseRewardClaims: len(tx.ScResults) != 2", txHash)
		return nil
	}
	rewardsIndex := 0
	if tx.SmartContractResults[0].Data == msgOKBase64 || tx.SmartContractResults[0].Data == msgOKHex {
		rewardsIndex = 1
	} else if tx.SmartContractResults[1].Data != msgOKBase64 && tx.SmartContractResults[1].Data != msgOKHex {
		log.Warn("Parser [tx_has: %s]: parseRewardClaims: can`t find OK msg", txHash)
		return nil
	}
	amount := node.ValueToEGLD(tx.SmartContractResults[rewardsIndex].Value)
	if tooMuchValue(amount) {
		log.Warn("Parser [tx_has: %s]: parseRewardClaims: too much value", txHash)
		return nil
	}
	d.rewards = append(d.rewards, dmodels.Reward{
		HypeblockID:     nonce,
		TxHash:          txHash,
		ReceiverAddress: tx.Sender,
		Amount:          amount,
		CreatedAt:       t,
	})
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    txHash,
		Type:      dmodels.ClaimRewardsEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     tx.Epoch,
		Amount:    amount,
		CreatedAt: t,
	})
	return nil
}

func (d *data) parseRewardDelegations(tx node.Tx, txHash string, nonce uint64, t time.Time) error {
	if !checkSCResults(tx.SmartContractResults, 2) {
		log.Warn("Parser: parseRewardDelegations: checkSCResults: false (tx: %s)", txHash)
		return nil
	}
	amount := tx.SmartContractResults[0].Value
	if len(tx.SmartContractResults[1].Data) == 0 {
		amount = tx.SmartContractResults[1].Value
	}
	amount = node.ValueToEGLD(amount)
	if tooMuchValue(amount) {
		log.Warn("Parser [tx_has: %s]: parseRewardDelegations: too much value", txHash)
		return nil
	}
	d.rewards = append(d.rewards, dmodels.Reward{
		HypeblockID:     nonce,
		TxHash:          txHash,
		ReceiverAddress: tx.Sender,
		Amount:          amount,
		CreatedAt:       t,
	})
	d.delegations = append(d.delegations, dmodels.Delegation{
		Delegator: tx.Sender,
		TxHash:    txHash,
		Validator: tx.Receiver,
		Amount:    amount,
		CreatedAt: t,
	})
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    txHash,
		Type:      dmodels.ReDelegateRewardsEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     tx.Epoch,
		Amount:    amount,
		CreatedAt: t,
	})
	return nil
}

func (d *data) parseUndelegations(tx node.Tx, txHash string, txType string, t time.Time) error {
	if !checkSCResults(tx.SmartContractResults, 2) {
		log.Warn("Parser: parseUndelegations: checkSCResults: false (tx: %s)", txHash)
		return nil
	}
	amountData := strings.TrimPrefix(txType, "unDelegate@")
	a, err := decimalFromHex(amountData)
	if err != nil {
		return fmt.Errorf("[tx: %s] decimalFromHex: %s", txHash, err.Error())
	}
	if tooMuchValue(a) {
		log.Warn("Parser [tx_has: %s]: parseUndelegations: too much value", txHash)
		return nil
	}
	d.delegations = append(d.delegations, dmodels.Delegation{
		Delegator: tx.Sender,
		TxHash:    txHash,
		Validator: tx.Receiver,
		Amount:    a.Neg(),
		CreatedAt: t,
	})
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    txHash,
		Type:      dmodels.UnDelegateStakeEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     tx.Epoch,
		Amount:    a.Neg(),
		CreatedAt: t,
	})
	return nil
}

func (d *data) parseStake(tx node.Tx, txHash string, t time.Time) error {
	if !checkSCResults(tx.SmartContractResults, 1) {
		log.Warn("Parser: parseStake: checkSCResults: false (tx: %s)", txHash)
		return nil
	}
	amount, err := decimal.NewFromString(tx.Value)
	if err != nil {
		return fmt.Errorf("decimal.NewFromString: %s", err.Error())
	}
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    txHash,
		Type:      dmodels.StakeStakeEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     tx.Epoch,
		Amount:    node.ValueToEGLD(amount),
		CreatedAt: t,
	})
	return nil
}

func (d *data) parseWithdraw(tx node.Tx, txHash string, t time.Time) error {
	findOK := false
	var amount decimal.Decimal
	for _, result := range tx.SmartContractResults {
		if result.Data == msgOKBase64 || result.Data == msgOKHex {
			findOK = true
		}
		if result.Receiver == tx.Sender {
			amount = node.ValueToEGLD(result.Value)
		}
	}
	if !findOK {
		return nil
	}
	if tooMuchValue(amount) {
		log.Warn("Parser [tx_has: %s]: parseWithdraw: too much value", txHash)
		return nil
	}
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    txHash,
		Type:      dmodels.WithdrawEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     tx.Epoch,
		Amount:    amount,
		CreatedAt: t,
	})
	return nil
}

func (d *data) parseUnstake(tx node.Tx, txType string, txHash string, t time.Time) error {
	if !checkSCResults(tx.SmartContractResults, 1) {
		log.Warn("Parser: parseUnstake: checkSCResults: false (tx: %s)", txHash)
		return nil
	}
	amountData := strings.TrimPrefix(txType, "unStake@")
	a, err := decimalFromHex(amountData)
	if err != nil {
		return fmt.Errorf("decimalFromHex: %s", err.Error())
	}
	if tooMuchValue(a) {
		log.Warn("Parser [tx_has: %s]: parseUnstake: too much value", txHash)
		return nil
	}
	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    txHash,
		Type:      dmodels.UnStakeEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     tx.Epoch,
		Amount:    a.Neg(),
		CreatedAt: t,
	})
	return nil
}

func (d *data) parseUnbond(tx node.Tx, txHash string, t time.Time) error {
	if len(tx.SmartContractResults) != 2 {
		log.Warn("Parser [tx_has: %s]: parseUnbond: len SmartContractResults != 2", txHash)
		return nil
	}
	okIndex := 1
	amountIndex := 0
	if base64.StdEncoding.EncodeToString([]byte(tx.SmartContractResults[1].Data)) == "delegation stake unbond" {
		okIndex = 0
		amountIndex = 1
	} else if base64.StdEncoding.EncodeToString([]byte(tx.SmartContractResults[0].Data)) != "delegation stake unbond" {
		log.Warn("Parser [tx_has: %s]: parseUnbond: can`t find `delegation stake unbond`", txHash)
		return nil
	}

	okStr := base64.StdEncoding.EncodeToString([]byte(tx.SmartContractResults[okIndex].Data))
	if !strings.Contains(okStr, "@ok") {
		log.Warn("Parser [tx_has: %s]: parseUnbond: bad OK", txHash)
		return nil
	}

	value := node.ValueToEGLD(tx.SmartContractResults[amountIndex].Value)
	if tooMuchValue(value) {
		log.Warn("Parser [tx_has: %s]: parseUnbond: too much value", txHash)
		return nil
	}

	d.stakeEvents = append(d.stakeEvents, dmodels.StakeEvent{
		TxHash:    txHash,
		Type:      dmodels.UnBondEventType,
		Validator: tx.Receiver,
		Delegator: tx.Sender,
		Epoch:     tx.Epoch,
		Amount:    value,
		CreatedAt: t,
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

	var dataset []data

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

		var singleData data
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

func checkSCResults(results []node.SmartContractResult, expectedLen int) bool {
	if len(results) != expectedLen {
		return false
	}
	switch expectedLen {
	case 1:
		return results[0].Data == msgOKBase64 || results[0].Data == msgOKHex
	case 2:
		okIndex := 0
		if len(results[0].Data) == 0 {
			okIndex = 1
		}
		return results[okIndex].Data == msgOKBase64 || results[okIndex].Data == msgOKHex
	}
	return false
}

func tooMuchValue(d decimal.Decimal) bool {
	return d.GreaterThanOrEqual(decimal.New(1, 18))
}
