package parser

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/config"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/shopspring/decimal"
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
		node      nodeAPI
		dao       dao.DAO
		fetcherCh chan uint64
		saverCh   chan data
		accounts  map[string]struct{}
		ctx       context.Context
		cancel    context.CancelFunc
		wg        *sync.WaitGroup
	}
	nodeAPI interface {
		GetTxByHash(hash string) (tx node.TxDetails, err error)
		GetTxsByMiniBlockHash(miniBlockHash string, offset, limit uint64) (txs []node.Tx, err error)
		GetMiniBlock(hash string) (miniBlock node.MiniBlock, err error)
		GetBlock(height uint64, shard uint64) (block node.Block, err error)
		GetBlockByHash(hash string, shard uint64) (block node.Block, err error)
		GetHyperBlock(height uint64) (hyperBlock node.HyperBlock, err error)
		GetMaxHeight(shardIndex uint64) (height uint64, err error)
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
		blockTimestamp := dmodels.NewTime(time.Unix(block.Block.Timestamp, 0))

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
			CreatedAt:       blockTimestamp,
		})

		for _, miniBlockInfo := range block.Block.Miniblocks {

			if _, ok := miniblocksRemain[miniBlockInfo.Hash]; ok {
				continue
			}
			miniblocksRemain[miniBlockInfo.Hash] = nil

			miniBlock, err := p.node.GetMiniBlock(miniBlockInfo.Hash)
			if err != nil {
				return d, fmt.Errorf("node.GetBlockByHash: %s", err.Error())
			}

			txs, err := p.node.GetTxsByMiniBlockHash(miniBlock.MiniBlockHash, 0, 1000) // todo check limit
			if err != nil {
				return d, fmt.Errorf("node.GetTxsByMiniBlockHash: %s", err.Error())
			}

			d.miniBlocks = append(d.miniBlocks, dmodels.MiniBlock{
				Hash:              miniBlock.MiniBlockHash,
				ReceiverBlockHash: miniBlock.ReceiverBlockHash,
				ReceiverShard:     miniBlock.ReceiverShard,
				SenderBlockHash:   miniBlock.SenderBlockHash,
				SenderShard:       miniBlock.SenderShard,
				Type:              miniBlock.Type,
				CreatedAt:         dmodels.NewTime(time.Unix(miniBlock.Timestamp, 0)),
			})

			for _, tx := range txs {
				switch strings.ToLower(tx.Status) {
				case dmodels.TxStatusPending:
					return d, fmt.Errorf("found pending tx: %s", tx.Txhash)
				case dmodels.TxStatusSuccess, dmodels.TxStatusFail:
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
					CreatedAt:     dmodels.NewTime(time.Unix(tx.Timestamp, 0)),
				})

				for _, r := range tx.ScResults {
					v, err := decimal.NewFromString(r.Value)
					if err != nil {
						return d, fmt.Errorf("decimal.NewFromString: %s", err.Error())
					}

					d.scResults = append(d.scResults, dmodels.SCResult{
						Hash:   r.Hash,
						TxHash: tx.Txhash,
						From:   r.Sender,
						To:     r.Receiver,
						Value:  v,
						Data:   r.Data,
					})
				}

				decodedBytes, err := base64.StdEncoding.DecodeString(tx.Data)
				if err != nil {
					return d, fmt.Errorf("base64.DecodeString: %s", err.Error())
				}
				if tx.Status == dmodels.TxStatusSuccess {
					txType := string(decodedBytes)
					switch txType {
					case "stake":
					case "reDelegateRewards":
					case "reStakeRewards": //	create stake, claimRewards
					case "delegate":
					case "claimRewards":
					case "unBond":
					default:
						if strings.Contains(txType, "unBond") {
							//fmt.Println(txType, tx.Txhash)
						}
						if strings.Contains(txType, "reStakeRewards") {
							//fmt.Println(txType, tx.Txhash)
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

		// todo
		// save via batch and timeout

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
			err = p.dao.UpdateParser(model)
			if err == nil {
				break
			}
			log.Error("Parser: dao.UpdateParser: %s", err.Error())
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
		accounts, err = p.dao.GetAccounts()
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
		addAccount(tx.Sender, tx.CreatedAt.Time)
		addAccount(tx.Receiver, tx.CreatedAt.Time)
	}

	for _, r := range d.scResults {
		addAccount(r.From, d.blocks[0].CreatedAt.Time)
		addAccount(r.To, d.blocks[0].CreatedAt.Time)
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
