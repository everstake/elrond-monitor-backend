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
	}
	block = smodels.Block{
	}
	return block, nil
}
