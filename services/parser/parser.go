package parser

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
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
	repeatDelay         = time.Second * 5
	parserTitle         = "elrond"
	fetcherChBuffer     = 5000
	saverChBuffer       = 5000
	metaChainShardIndex = 4294967295
	precision           = 18
)

var shardIndexes = map[uint64]ShardIndex{
	0:                   0,
	1:                   1,
	2:                   2,
	metaChainShardIndex: metaChainShardIndex,
}

var precisionDiv = decimal.New(1, precision)

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
	}
	data struct {
		height       uint64
		blocks       []dmodels.Block
		miniBlocks   []dmodels.MiniBlock
		transactions []dmodels.Transaction
		scResults    []dmodels.SCResult
		accounts     []dmodels.Account
		stakes       []dmodels.Stake
		delegations  []dmodels.Delegation
		rewards      []dmodels.Reward
	}
	ShardIndex uint64
)

func NewParser(cfg config.Config, d dao.DAO) *Parser {
	ctx, cancel := context.WithCancel(context.Background())
	return &Parser{
		cfg:       cfg,
		dao:       d,
		node:      node.NewAPI(cfg.Parser.Node),
		fetcherCh: make(chan uint64, fetcherChBuffer),
		saverCh:   make(chan data, saverChBuffer),
		accounts:  make(map[string]struct{}),
		ctx:       ctx,
		cancel:    cancel,
		wg:        &sync.WaitGroup{},
	}
}

func (p *Parser) Run() error {
	model, err := p.dao.GetParser(parserTitle)
	if err != nil {
		return fmt.Errorf("parser not found")
	}
	for i := uint64(0); i < p.cfg.Parser.Fetchers; i++ {
		go p.runFetcher()
	}

	go p.saving()
	for {
		latestBlock, err := p.node.GetMaxHeight(metaChainShardIndex)
		if err != nil {
			log.Error("Parser: node.GetMaxHeight: %s", err.Error())
			continue
		}
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

	hyperBlocks := make([]node.Block, 0, len(shardIndexes))
	metaChainBlock, err := p.node.GetBlockByHash(hyperBlock.HyperBlock.Hash, metaChainShardIndex)
	if err != nil {
		return d, fmt.Errorf("api.GetBlockByHash: %s", err.Error())
	}
	hyperBlocks = append(hyperBlocks, metaChainBlock)
	for _, ShardBlockInfo := range hyperBlock.HyperBlock.Shardblocks {
		block, err := p.node.GetBlockByHash(ShardBlockInfo.Hash, ShardBlockInfo.Shard)
		if err != nil {
			return d, fmt.Errorf("api.GetBlockByHash: %s", err.Error())
		}
		hyperBlocks = append(hyperBlocks, block)
	}

	miniblocksRemain := make(map[string]interface{})

	for _, block := range hyperBlocks {
		d.blocks = append(d.blocks, dmodels.Block{
			AccumulatedFees: decimal.New(block.Block.AccumulatedFees, 0).Div(precisionDiv),
			DeveloperFees:   decimal.New(block.Block.DeveloperFees, 0).Div(precisionDiv),
			Hash:            block.Block.Hash,
			Nonce:           block.Block.Nonce,
			Round:           block.Block.Round,
			Shard:           block.Block.Shard,
			NumTxs:          block.Block.NumTxs,
			Epoch:           block.Block.Epoch,
			Status:          block.Block.Status,
			PrevBlockHash:   block.Block.PrevBlockHash,
			CreatedAt:       time.Unix(block.Block.Timestamp, 0),
		})

		for _, miniBlockInfo := range block.Block.Miniblocks {
			// ignore same miniblocks
			if _, ok := miniblocksRemain[miniBlockInfo.Hash]; ok {
				continue
			}
			miniblocksRemain[miniBlockInfo.Hash] = nil

			miniBlock, err := p.node.GetMiniBlock(miniBlockInfo.Hash)
			if err != nil {
				return d, fmt.Errorf("node.GetBlockByHash: %s", err.Error())
			}

			txs, err := p.node.GetTxsByMiniBlockHash(miniBlock.MiniBlockHash, 0, 1000)
			if err != nil {
				return d, fmt.Errorf("node.GetTxsByMiniBlockHash: %s", err.Error())
			}

			// check number of txs from mini block
			if len(txs) == 1000 {
				return d, fmt.Errorf("the maximum number of transactions has been reached")
			}

			d.miniBlocks = append(d.miniBlocks, dmodels.MiniBlock{
				Hash:              miniBlock.MiniBlockHash,
				ReceiverBlockHash: miniBlock.ReceiverBlockHash,
				ReceiverShard:     miniBlock.ReceiverShard,
				SenderBlockHash:   miniBlock.SenderBlockHash,
				SenderShard:       miniBlock.SenderShard,
				Type:              miniBlock.Type,
				CreatedAt:         time.Unix(miniBlock.Timestamp, 0),
			})

			for _, tx := range txs {
				switch strings.ToLower(tx.Status) {
				case dmodels.TxStatusPending:
					return d, fmt.Errorf("found pending tx: %s", tx.Txhash)
				case dmodels.TxStatusSuccess, dmodels.TxStatusFail, dmodels.TxStatusInvalid:
				default:
					return d, fmt.Errorf("unknown tx status: %s", tx.Status)
				}

				amount, err := decimal.NewFromString(tx.Value)
				if err != nil {
					return d, fmt.Errorf("decimal.NewFromString: %s", err.Error())
				}

				d.transactions = append(d.transactions, dmodels.Transaction{
					Hash:          tx.Txhash,
					Status:        tx.Status,
					MiniBlockHash: tx.MiniBlockHash,
					Value:         amount.Div(precisionDiv),
					Fee:           decimal.New(tx.Fee, 0).Div(precisionDiv),
					Sender:        tx.Sender,
					SenderShard:   tx.SenderShard,
					Receiver:      tx.Receiver,
					ReceiverShard: tx.ReceiverShard,
					GasPrice:      tx.GasPrice,
					GasUsed:       tx.GasUsed,
					Nonce:         tx.Nonce,
					Data:          tx.Data,
					CreatedAt:     time.Unix(tx.Timestamp, 0),
				})

				for _, r := range tx.ScResults {
					v, err := decimal.NewFromString(r.Value)
					if err != nil {
						return d, fmt.Errorf("decimal.NewFromString: %s", err.Error())
					}

					d.scResults = append(d.scResults, dmodels.SCResult{
						Hash:    r.Hash,
						TxHash:  tx.Txhash,
						From:    r.Sender,
						To:      r.Receiver,
						Value:   v,
						Data:    r.Data,
						Message: r.ReturnMessage,
					})
				}

				decodedBytes, err := base64.StdEncoding.DecodeString(tx.Data)
				if err != nil {
					return d, fmt.Errorf("base64.DecodeString: %s", err.Error())
				}

				/*
				     mixed sorting sc_results

					 examples:
						delegate - 44b7729c15b4ae36e56d742ed81d2510a347033f73e9c5db2d117917c4996a13
						unDelegate - 596155353284baf98a5b9a539ab941898ac58c36f8ce54c642af5b264aac8338
						reDelegateRewards 10c8d8f23973ff3a00ae86fc0f4b6ae2a70105d46ecdb702cabfdd99e27363d4
						claimRewards - 743c6d62f0037d876e3f41284e8af2595ce5d63bc0b1bf08281b36f182c3bb83
						unStake - d10adba96c063a55c1b369094073c41ae57286fab53c25cf8489d9c4c4ffbb18
						stake - 7c7db6eea2b3f2aef9875f91300e149b5b379d86456110eca57545f3aa087886
						unBond - c6c32820df1b44af828121d52207904c53230e9dd3a794e90e714915a68806d4

				*/
				if tx.Status == dmodels.TxStatusSuccess {
					txType := string(decodedBytes)
					switch txType {
					case "withdraw":
						// todo research
					case "stake":
						d.stakes = append(d.stakes, dmodels.Stake{
							ID:        tx.Txhash,
							TxHash:    tx.Txhash,
							Delegator: tx.Sender,
							Validator: tx.Receiver,
							Amount:    amount.Div(precisionDiv),
							CreatedAt: time.Unix(tx.Timestamp, 0),
						})
					case "reDelegateRewards": // create delegation + claimRewards
						if len(tx.ScResults) != 2 {
							continue
						}
						if tx.Sender != tx.ScResults[1].Receiver {
							fmt.Println("reDelegateRewards different delegators") // todo delete it
							continue
						}
						a, err := decimal.NewFromString(tx.ScResults[0].Value)
						if err != nil {
							return d, fmt.Errorf("decimal.NewFromString: %s", err.Error())
						}
						d.rewards = append(d.rewards, dmodels.Reward{
							ID:              tx.Txhash,
							HypeblockID:     nonce,
							TxHash:          tx.Txhash,
							ReceiverAddress: tx.Sender,
							Amount:          a.Div(precisionDiv),
							CreatedAt:       time.Unix(tx.Timestamp, 0),
						})
						d.delegations = append(d.delegations, dmodels.Delegation{
							ID:        tx.Txhash,
							Delegator: tx.Sender,
							TxHash:    tx.Txhash,
							Validator: tx.Receiver,
							Amount:    a.Div(precisionDiv),
							CreatedAt: time.Unix(tx.Timestamp, 0),
						})
					case "reStakeRewards": // create stake + claimRewards (check existence of reStakeRewards tx)
						fmt.Println(txType, tx.Txhash)
					case "delegate":
						if len(tx.ScResults) != 2 {
							continue
						}
						if tx.Sender != tx.ScResults[1].Receiver {
							fmt.Println("delegate different delegators") // todo delete it
							continue
						}
						d.delegations = append(d.delegations, dmodels.Delegation{
							ID:        tx.Txhash,
							Delegator: tx.Sender,
							TxHash:    tx.Txhash,
							Validator: tx.Receiver,
							Amount:    amount.Div(precisionDiv),
							CreatedAt: time.Unix(tx.Timestamp, 0),
						})
					case "claimRewards":
						if len(tx.ScResults) != 2 {
							continue
						}
						if tx.Sender != tx.ScResults[0].Receiver {
							fmt.Println("claimRewards different delegators") // todo delete it
							continue
						}
						reward, err := decimal.NewFromString(tx.ScResults[1].Value)
						if err != nil {
							return d, fmt.Errorf("decimal.NewFromString: %s", err.Error())
						}
						d.rewards = append(d.rewards, dmodels.Reward{
							ID:              tx.Txhash,
							HypeblockID:     nonce,
							TxHash:          tx.Txhash,
							ReceiverAddress: tx.Sender,
							Amount:          reward.Div(precisionDiv),
							CreatedAt:       time.Unix(tx.Timestamp, 0),
						})
					case "unBond":
						fmt.Println(txType, tx.Txhash)
					default:
						if strings.Contains(txType, "unBond") {
							fmt.Println(txType, 2, tx.Txhash)
						}
						if strings.Contains(txType, "relayedTx") {
							//fmt.Println(txType, tx.Txhash)
						}
						if strings.Contains(txType, "unDelegate") {
							amountData := strings.TrimLeft(txType, "unDelegate@")
							a := (&big.Int{}).SetBytes([]byte(amountData))
							d.delegations = append(d.delegations, dmodels.Delegation{
								ID:        tx.Txhash,
								Delegator: tx.Sender,
								TxHash:    tx.Txhash,
								Validator: tx.Receiver,
								Amount:    decimal.NewFromBigInt(a, 0).Div(precisionDiv).Neg(),
								CreatedAt: time.Unix(tx.Timestamp, 0),
							})
						}
						if strings.Contains(txType, "unStake") {
							amountData := strings.TrimLeft(txType, "unStake@")
							a := (&big.Int{}).SetBytes([]byte(amountData))
							d.stakes = append(d.stakes, dmodels.Stake{
								ID:        tx.Txhash,
								TxHash:    tx.Txhash,
								Delegator: tx.Sender,
								Validator: tx.Receiver,
								Amount:    decimal.NewFromBigInt(a, 0).Div(precisionDiv).Neg(),
								CreatedAt: time.Unix(tx.Timestamp, 0),
							})
						}
						if strings.Contains(txType, "reStakeRewards") { // check existence of reStakeRewards tx
							fmt.Println(txType, tx.Txhash)
						}
					}
				}
			}
		}

	}

	return d, nil
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
	p.setAccounts()

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
			singleData.blocks = append(singleData.blocks, item.blocks...)
			singleData.miniBlocks = append(singleData.miniBlocks, item.miniBlocks...)
			singleData.transactions = append(singleData.transactions, item.transactions...)
			singleData.scResults = append(singleData.scResults, item.scResults...)
			singleData.accounts = append(singleData.accounts, item.accounts...)
		}
		p.wg.Add(1)
		var err error
		for {
			err = p.dao.CreateBlocks(singleData.blocks)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateBlocks: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateMiniBlocks(singleData.miniBlocks)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateMiniBlocks: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateTransactions(singleData.transactions)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateTransactions: %s", err.Error())
			<-time.After(repeatDelay)
		}
		for {
			err = p.dao.CreateSCResults(singleData.scResults)
			if err == nil {
				break
			}
			log.Error("Parser: dao.CreateSCResults: %s", err.Error())
			<-time.After(repeatDelay)
		}

		p.saveNewAccounts(singleData)
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

func (p *Parser) setAccounts() {
	var accounts []dmodels.Account
	var err error
	for {
		accounts, err = p.dao.GetAccounts(filters.Accounts{})
		if err != nil {
			log.Error("Parser: setAccounts: dao.GetAccounts: %s", err.Error())
			<-time.After(repeatDelay)
			continue
		}
		break
	}
	for _, account := range accounts {
		p.accounts[account.Address] = struct{}{}
	}
}

func (p *Parser) saveNewAccounts(d data) {
	var newAccounts []dmodels.Account
	addAccount := func(acc string, tm time.Time) {
		_, ok := p.accounts[acc]
		if !ok {
			p.accounts[acc] = struct{}{}
			newAccounts = append(newAccounts, dmodels.Account{
				Address:   acc,
				CreatedAt: tm,
			})
		}
	}

	for _, tx := range d.transactions {
		addAccount(tx.Sender, tx.CreatedAt)
		addAccount(tx.Receiver, tx.CreatedAt)
	}

	for _, r := range d.scResults {
		addAccount(r.From, d.blocks[0].CreatedAt)
		addAccount(r.To, d.blocks[0].CreatedAt)
	}

	for {
		err := p.dao.CreateAccounts(newAccounts)
		if err == nil {
			break
		}
		log.Error("Parser: dao.CreateAccounts: %s", err.Error())
		<-time.After(repeatDelay)
	}
}
