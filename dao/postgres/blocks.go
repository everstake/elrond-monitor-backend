package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
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
	conflictCondition := `
		ON CONFLICT (mlk_hash) DO UPDATE SET 
			mlk_receiver_block_hash = CASE WHEN miniblocks.mlk_receiver_block_hash = '' THEN EXCLUDED.mlk_receiver_block_hash ELSE miniblocks.mlk_receiver_block_hash END,
			mlk_sender_block_hash = CASE WHEN miniblocks.mlk_sender_block_hash = '' THEN EXCLUDED.mlk_sender_block_hash ELSE miniblocks.mlk_sender_block_hash END
	`
	q = q.Suffix(conflictCondition)
	_, err := db.insert(q)
	return err
}

func (db Postgres) GetBlocks(filter filters.Blocks) (blocks []dmodels.Block, err error) {
	q := squirrel.Select("*").From(dmodels.BlocksTable).OrderBy("blk_created_at desc")
	if len(filter.Shard) != 0 {
		q = q.Where(squirrel.Eq{"blk_shard": filter.Shard})
	}
	if filter.Nonce != 0 {
		q = q.Where(squirrel.Eq{"blk_nonce": filter.Nonce})
	}
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset() != 0 {
		q = q.Offset(filter.Offset())
	}
	err = db.find(&blocks, q)
	return blocks, err
}

func (db Postgres) GetBlock(hash string) (block dmodels.Block, err error) {
	q := squirrel.Select("*").From(dmodels.BlocksTable).Where(squirrel.Eq{"blk_hash": hash})
	err = db.first(&block, q)
	return block, err
}

func (db Postgres) GetBlocksTotal(filter filters.Blocks) (total uint64, err error) {
	q := squirrel.Select("count(*) as total").From(dmodels.BlocksTable)
	if len(filter.Shard) != 0 {
		q = q.Where(squirrel.Eq{"blk_shard": filter.Shard})
	}
	if filter.Nonce != 0 {
		q = q.Where(squirrel.Eq{"blk_nonce": filter.Nonce})
	}
	err = db.first(&total, q)
	return total, err
}

func (db Postgres) GetMiniBlocks(filter filters.MiniBlocks) (blocks []dmodels.MiniBlock, err error) {
	q := squirrel.Select("*").From(dmodels.MiniBlocksTable)
	if filter.ParentBlockHash != "" {
		q = q.Where(squirrel.Or{squirrel.Eq{"mlk_sender_block_hash": filter.ParentBlockHash}, squirrel.Eq{"mlk_receiver_block_hash": filter.ParentBlockHash}})
	}
	err = db.find(&blocks, q)
	return blocks, err
}

func (db Postgres) GetMiniBlock(hash string) (block dmodels.MiniBlock, err error) {
	q := squirrel.Select("*").From(dmodels.MiniBlocksTable).Where(squirrel.Eq{"mlk_hash": hash})
	err = db.first(&block, q)
	return block, err
}

func (db Postgres) GetMiniBlocksTotal(filter filters.MiniBlocks) (total uint64, err error) {
	q := squirrel.Select("count(*) as total").From(dmodels.BlocksTable)
	if filter.ParentBlockHash != "" {
		q = q.Where(squirrel.Or{squirrel.Eq{"mlk_sender_block_hash": filter.ParentBlockHash}, squirrel.Eq{"mlk_receiver_block_hash": filter.ParentBlockHash}})
	}
	err = db.first(&total, q)
	return total, err
}
