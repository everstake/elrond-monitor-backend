package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
)

func (db Postgres) CreateBlocks(blocks []dmodels.Block) error {
	if len(blocks) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.BlocksTable).Columns(
		"blk_hash",
		"blk_nonce",
		"blk_round",
		"blk_shard",
		"blk_num_txs",
		"blk_epoch",
		"blk_status",
		"blk_prev_block_hash",
		"blk_accumulated_fees",
		"blk_developer_fees",
		"blk_created_at",
	)
	for _, b := range blocks {
		if b.Hash == "" {
			return fmt.Errorf("field Hash is empty")
		}
		if b.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			b.Hash,
			b.Nonce,
			b.Round,
			b.Shard,
			b.NumTxs,
			b.Epoch,
			b.Status,
			b.PrevBlockHash,
			b.AccumulatedFees,
			b.DeveloperFees,
			b.CreatedAt,
		)
	}
	q = q.Suffix("ON CONFLICT (blk_hash) DO NOTHING")
	_, err := db.insert(q)
	return err
}

func (db Postgres) CreateMiniBlocks(blocks []dmodels.MiniBlock) error {
	if len(blocks) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.MiniBlocksTable).Columns(
		"mlk_hash",
		"mlk_receiver_block_hash",
		"mlk_receiver_shard",
		"mlk_sender_block_hash",
		"mlk_sender_shard",
		"mlk_type",
		"mlk_created_at",
	)
	for _, b := range blocks {
		if b.Hash == "" {
			return fmt.Errorf("field Hash is empty")
		}
		if b.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			b.Hash,
			b.ReceiverBlockHash,
			b.ReceiverShard,
			b.SenderBlockHash,
			b.SenderShard,
			b.Type,
			b.CreatedAt,
		)
	}
	q = q.Suffix("ON CONFLICT (mlk_hash) DO NOTHING")
	_, err := db.insert(q)
	return err
}
