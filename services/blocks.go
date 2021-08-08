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
	esBlock, err := s.es.GetBlock(hash)
	if err != nil {
		return block, fmt.Errorf("es.GetBlock: %s", err.Error())
	}
	esValidatorsKeys, err := s.es.ValidatorsKeys(dBlock.Shard, dBlock.Epoch)
	if err != nil {
		return block, fmt.Errorf("es.ValidatorsKeys: %s", err.Error())
	}
	validatorsKeys := make([]string, len(esValidatorsKeys.PublicKeys))
	for i, key := range esBlock.Validators {
		validatorsKeys[i] = validatorKeyByIndex(esValidatorsKeys.PublicKeys, key)
	}
	block = smodels.Block{
		Hash:                  dBlock.Hash,
		Nonce:                 dBlock.Nonce,
		Shard:                 dBlock.Shard,
		Epoch:                 dBlock.Epoch,
		TxCount:               dBlock.NumTxs,
		Size:                  esBlock.Size,
		Proposer:              validatorKeyByIndex(esValidatorsKeys.PublicKeys, esBlock.Proposer),
		Miniblocks:            miniBlocksHashes,
		NotarizedBlocksHashes: esBlock.NotarizedBlocksHashes,
		Validators:            validatorsKeys,
		PubKeyBitmap:          esBlock.PubKeyBitmap,
		StateRootHash:         esBlock.StateRootHash,
		PrevHash:              esBlock.PrevHash,
		Timestamp:             smodels.NewTime(dBlock.CreatedAt),
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
	dBlock, err := s.dao.GetMiniBlock(hash)
	if err != nil {
		return block, fmt.Errorf("dao.GetMiniBlock: %s", err.Error())
	}
	dTxs, err := s.dao.GetTransactions(filters.Transactions{MiniBlock: hash})
	if err != nil {
		return block, fmt.Errorf("dao.GetTransactions: %s", err.Error())
	}
	txs := make([]smodels.Tx, len(dTxs))
	for i, tx := range dTxs {
		txs[i] = smodels.Tx{
			Hash:          tx.Hash,
			Status:        tx.Status,
			From:          tx.Sender,
			To:            tx.Receiver,
			Value:         tx.Value,
			MiniblockHash: tx.MiniBlockHash,
			ShardFrom:     tx.SenderShard,
			ShardTo:       tx.ReceiverShard,
			Type:          "", // todo
			Timestamp:     smodels.NewTime(tx.CreatedAt),
		}
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
