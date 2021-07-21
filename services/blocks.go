package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/smodels"
)

func (s *ServiceFacade) GetBlock(hash string) (block smodels.Block, err error) {
	dBlock, err := s.dao.GetBlock(hash)
	if err != nil {
		return block, fmt.Errorf("dao.GetBlock: %s", err.Error())
	}
	dMiniBlocks, err := s.dao.GetMiniBlocks(filters.MiniBlocks{ParentBlockHash: dBlock.Hash})
	if err != nil {
		return block, fmt.Errorf("dao.GetMiniBlocks: %s", err.Error())
	}
	miniBlocksHashes := make([]string, len(dMiniBlocks))
	miniBlocks := make([]smodels.Miniblock, len(dMiniBlocks))
	for i, b := range dMiniBlocks {
		miniBlocks[i] = smodels.Miniblock{
			Hash:          b.Hash,
			ShardFrom:     b.SenderShard,
			ShardTo:       b.ReceiverShard,
			BlockSender:   b.SenderBlockHash,
			BlockReceiver: b.ReceiverBlockHash,
			Type:          b.Type,
			Timestamp:     smodels.NewTime(b.CreatedAt),
		}
		miniBlocksHashes[i] = b.Hash
	}
	//extraData, err := s.node.GetExtraDataBlock(dBlock.Hash)
	//if err != nil {
	//	return block, fmt.Errorf("node.GetExtraDataBlock: %s", err.Error())
	//}
	block = smodels.Block{
		Hash:       dBlock.Hash,
		Nonce:      dBlock.Nonce,
		Shard:      dBlock.Shard,
		Epoch:      dBlock.Epoch,
		TxCount:    dBlock.NumTxs,
		//Size:       extraData.Size,
		//Proposer:   extraData.Proposer,
		Miniblocks: miniBlocksHashes,
		Timestamp:  smodels.NewTime(dBlock.CreatedAt),
	}
	return block, nil
}

func (s *ServiceFacade) GetBlocks(filter filters.Blocks) (items smodels.Pagination, err error) {
	dBlocks, err := s.dao.GetBlocks(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetBlocks: %s", err.Error())
	}
	blocks := make([]smodels.Block, len(dBlocks))
	for i, b := range dBlocks {
		blocks[i] = smodels.Block{
			Hash:      b.Hash,
			Nonce:     b.Nonce,
			Shard:     b.Shard,
			Epoch:     b.Epoch,
			TxCount:   b.NumTxs,
			Timestamp: smodels.NewTime(b.CreatedAt),
		}
	}
	total, err := s.dao.GetBlocksTotal(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetBlocksTotal: %s", err.Error())
	}
	return smodels.Pagination{
		Items: blocks,
		Count: total,
	}, nil
}

func (s *ServiceFacade) GetMiniBlock(hash string) (block smodels.Miniblock, err error) {
	dBlock, err := s.dao.GetMiniBlock(hash)
	if err != nil {
		return block, fmt.Errorf("dao.GetMiniBlock: %s", err.Error())
	}
	dTxs, err := s.dao.GetTransactions(filters.Transactions{MiniBlock: hash})
	if err != nil {
		return block, fmt.Errorf("dao.GetTransactions: %s", err.Error())
	}
	txs := make([]smodels.Tx, len(dTxs))
	for _, tx := range dTxs {
		txs = append(txs, smodels.Tx{
			Hash:          tx.Hash,
			Status:        tx.Status,
			From:          tx.Sender,
			To:            tx.Receiver,
			Value:         tx.Value,
			Fee:           tx.Fee,
			GasUsed:       tx.GasUsed,
			MiniblockHash: tx.MiniBlockHash,
			ShardFrom:     tx.SenderShard,
			ShardTo:       tx.ReceiverShard,
			Type:          "", // todo
			Timestamp:     smodels.NewTime(tx.CreatedAt),
		})
	}
	return smodels.Miniblock{
		Hash:          dBlock.Hash,
		ShardFrom:     dBlock.SenderShard,
		ShardTo:       dBlock.ReceiverShard,
		BlockSender:   dBlock.SenderBlockHash,
		BlockReceiver: dBlock.ReceiverBlockHash,
		Type:          dBlock.Type,
		Txs:           txs,
		Timestamp:     smodels.NewTime(dBlock.CreatedAt),
	}, nil
}
