package watcher

import (
	"fmt"
	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/everstake/elrond-monitor-backend/api/ws"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"time"
)

const interval = time.Second * 3

type Watcher struct {
	dao           dao.DAO
	lastBlockTime int64
	lastTxTime    int64
	stop          chan struct{}
	ws            ws.WS
}

func NewWatcher(d dao.DAO, w ws.WS) *Watcher {
	return &Watcher{
		dao:  d,
		ws:   w,
		stop: make(chan struct{}),
	}
}

func (w *Watcher) Run() (err error) {
	blocks, err := w.dao.GetBlocks(filters.Blocks{
		Pagination: filters.Pagination{Limit: 1},
	})
	if err != nil {
		return fmt.Errorf("dao.GetBlocks: %s", err.Error())
	}
	if len(blocks) == 0 {
		return fmt.Errorf("total blocks is zero")
	}
	w.lastBlockTime = int64(blocks[0].Timestamp)

	txs, err := w.dao.GetTransactions(filters.Transactions{
		Pagination: filters.Pagination{Limit: 1},
	})
	if err != nil {
		return fmt.Errorf("dao.GetTransactions: %s", err.Error())
	}
	if len(blocks) == 0 {
		return fmt.Errorf("total txs is zero")
	}
	w.lastTxTime = int64(txs[0].Timestamp)

	for {
		select {
		case <-w.stop:
			return nil
		case <-time.After(interval):
			// blocks
			blocks, err = w.dao.GetBlocks(filters.Blocks{
				Pagination: filters.Pagination{Limit: 10},
			})
			if err != nil {
				log.Warn("Watcher: dao.GetBlocks: %s", err.Error())
				continue
			}
			var newBlocks []data.Block
			maxBlockTime := w.lastBlockTime
			for _, block := range blocks {
				t := int64(block.Timestamp)
				if t > w.lastBlockTime {
					newBlocks = append(newBlocks, block)
				}
				if t > maxBlockTime {
					maxBlockTime = t
				}
			}
			if len(newBlocks) > 0 {
				w.ws.Broadcast(ws.Broadcast{
					Channel: ws.BlocksChannel,
					Data:    newBlocks,
				})
				w.lastBlockTime = maxBlockTime
			}

			// transactions
			txs, err = w.dao.GetTransactions(filters.Transactions{
				Pagination: filters.Pagination{Limit: 10},
			})
			if err != nil {
				log.Warn("Watcher: dao.GetTransactions: %s", err.Error())
				continue
			}
			var newTxs []data.Transaction
			maxTxTime := w.lastTxTime
			for _, tx := range txs {
				t := int64(tx.Timestamp)
				if t > w.lastTxTime {
					newTxs = append(newTxs, tx)
				}
				if t > maxTxTime {
					maxTxTime = t
				}
			}
			if len(newTxs) > 0 {
				w.ws.Broadcast(ws.Broadcast{
					Channel: ws.TransactionsChannel,
					Data:    newTxs,
				})
				w.lastTxTime = maxTxTime
			}
		}
	}
}

func (w *Watcher) Stop() error {
	w.stop <- struct{}{}
	return nil
}

func (w *Watcher) Title() string {
	return "Watcher"
}
