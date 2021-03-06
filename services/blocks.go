package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/derrors"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

func (s *ServiceFacade) GetBlock(hash string) (block smodels.Block, err error) {
	dBlock, err := s.dao.GetBlock(hash)
	if err != nil {
		if err == derrors.NotFound {
			return block, smodels.Error{
				Err:      err.Error(),
				Msg:      "block not found",
				HttpCode: http.StatusNotFound,
			}
		}
		return block, fmt.Errorf("dao.GetBlock: %s", err.Error())
	}
	esValidatorsKeys, err := s.dao.ValidatorsKeys(uint64(dBlock.ShardID), uint64(dBlock.Epoch))
	if err != nil {
		return block, fmt.Errorf("es.ValidatorsKeys: %s", err.Error())
	}
	var validatorsKeys []string
	for _, key := range dBlock.Validators {
		val := validatorKeyByIndex(esValidatorsKeys.PublicKeys, key)
		if val != "" {
			validatorsKeys = append(validatorsKeys, val)
		}
	}
	block = smodels.Block{
		Hash:                  hash,
		Nonce:                 dBlock.Nonce,
		Shard:                 uint64(dBlock.ShardID),
		Epoch:                 uint64(dBlock.Epoch),
		TxCount:               uint64(dBlock.TxCount),
		Size:                  dBlock.Size,
		Proposer:              validatorKeyByIndex(esValidatorsKeys.PublicKeys, dBlock.Proposer),
		Miniblocks:            dBlock.MiniBlocksHashes,
		NotarizedBlocksHashes: dBlock.NotarizedBlocksHashes,
		Validators:            validatorsKeys,
		PubKeyBitmap:          dBlock.PubKeyBitmap,
		StateRootHash:         dBlock.StateRootHash,
		PrevHash:              dBlock.PrevHash,
		Timestamp:             smodels.NewTime(time.Unix(int64(dBlock.Timestamp), 0)),
	}
	return block, nil
}

func validatorKeyByIndex(keys []string, index uint64) string {
	if uint64(len(keys)) <= index {
		return ""
	}
	return keys[index]
}

func (s *ServiceFacade) GetBlocks(filter filters.Blocks) (items smodels.Pagination, err error) {
	dBlocks, err := s.dao.GetBlocks(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetBlocks: %s", err.Error())
	}
	blocks := make([]smodels.Block, len(dBlocks))
	for i, b := range dBlocks {
		blocks[i] = smodels.Block{
			Hash:                  b.Hash,
			Nonce:                 b.Nonce,
			Shard:                 uint64(b.ShardID),
			Epoch:                 uint64(b.Epoch),
			TxCount:               uint64(b.TxCount),
			Size:                  b.Size,
			Miniblocks:            b.MiniBlocksHashes,
			NotarizedBlocksHashes: b.NotarizedBlocksHashes,
			PubKeyBitmap:          b.PubKeyBitmap,
			StateRootHash:         b.StateRootHash,
			PrevHash:              b.PrevHash,
			Timestamp:             smodels.NewTime(time.Unix(int64(b.Timestamp), 0)),
		}
	}
	total, err := s.dao.GetBlocksCount(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetBlocksCount: %s", err.Error())
	}
	return smodels.Pagination{
		Items: blocks,
		Count: total,
	}, nil
}

func (s *ServiceFacade) GetBlockByNonce(shard uint64, nonce uint64) (block smodels.Block, err error) {
	dBlocks, err := s.dao.GetBlocks(filters.Blocks{
		Shard: []uint64{shard},
		Nonce: nonce,
	})
	if err != nil {
		return block, fmt.Errorf("dao.GetBlocks: %s", err.Error())
	}
	if len(dBlocks) == 0 {
		errMsg := fmt.Sprintf("not found shard: %d, nonce: %d", shard, nonce)
		return block, smodels.Error{Err: errMsg, Msg: errMsg, HttpCode: 404}
	}
	return s.GetBlock(dBlocks[0].Hash)
}

func (s *ServiceFacade) GetMiniBlock(hash string) (block smodels.Miniblock, err error) {
	dBlock, err := s.dao.GetMiniblock(hash)
	if err != nil {
		if err == derrors.NotFound {
			return block, smodels.Error{
				Err:      err.Error(),
				Msg:      "miniblock not found",
				HttpCode: http.StatusNotFound,
			}
		}
		return block, fmt.Errorf("dao.GetMiniblock: %s", err.Error())
	}
	dTxs, err := s.dao.GetTransactions(filters.Transactions{MiniBlock: hash})
	if err != nil {
		return block, fmt.Errorf("dao.GetTransactions: %s", err.Error())
	}
	txs := make([]smodels.Tx, len(dTxs))
	for i, tx := range dTxs {
		val, _ := decimal.NewFromString(tx.Value)
		fee, _ := decimal.NewFromString(tx.Fee)
		txs[i] = smodels.Tx{
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
			Timestamp:     smodels.NewTime(time.Unix(int64(tx.Timestamp), 0)),
		}
	}
	return smodels.Miniblock{
		Hash:          hash,
		ShardFrom:     uint64(dBlock.SenderShardID),
		ShardTo:       uint64(dBlock.ReceiverShardID),
		BlockSender:   dBlock.SenderBlockHash,
		BlockReceiver: dBlock.ReceiverBlockHash,
		Type:          dBlock.Type,
		Txs:           txs,
		Timestamp:     smodels.NewTime(time.Unix(int64(dBlock.Timestamp), 0)),
	}, nil
}
