package watcher

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/api/ws"
	"github.com/everstake/elrond-monitor-backend/dao"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/log"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
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
			var newBlocks []smodels.Block
			maxBlockTime := w.lastBlockTime
			for _, block := range blocks {
				t := int64(block.Timestamp)
				if t > w.lastBlockTime {
					newBlocks = append(newBlocks, smodels.Block{
						Hash:          block.Hash,
						Nonce:         block.Nonce,
						Shard:         uint64(block.ShardID),
						Epoch:         uint64(block.Epoch),
						TxCount:       uint64(block.TxCount),
						Size:          block.Size,
						Miniblocks:    block.MiniBlocksHashes,
						PubKeyBitmap:  block.PubKeyBitmap,
						StateRootHash: block.StateRootHash,
						PrevHash:      block.PrevHash,
						Timestamp:     smodels.NewTime(time.Unix(int64(block.Timestamp), 0)),
					})
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
			var newTxs []smodels.Tx
			maxTxTime := w.lastTxTime
			for _, tx := range txs {
				t := int64(tx.Timestamp)
				if t > w.lastTxTime {
					val, _ := decimal.NewFromString(tx.Value)
					fee, _ := decimal.NewFromString(tx.Fee)
					newTxs = append(newTxs, smodels.Tx{
						Hash:          tx.Hash,
						Status:        tx.Status,
						From:          tx.Sender,
						To:            tx.Receiver,
						Value:         node.ValueToEGLD(val),
						Fee:           node.ValueToEGLD(fee),
						GasUsed:       tx.GasUsed,
						GasPrice:      tx.GasPrice,
						MiniblockHash: tx.MBHash,
						ShardFrom:     uint64(tx.SenderShard),
						ShardTo:       uint64(tx.ReceiverShard),
						Signature:     tx.Signature,
						Data:          string(tx.Data),
						Timestamp:     smodels.NewTime(time.Unix(int64(tx.Timestamp), 0)),
					})
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
