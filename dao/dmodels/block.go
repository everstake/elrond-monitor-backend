package dmodels

import "github.com/shopspring/decimal"

const (
	BlocksTable     = "blocks"
	MiniBlocksTable = "miniblocks" // type == TxBlock
)

type (
	Block struct {
		Hash            string          `db:"blk_hash" json:"hash"`
		Nonce           uint64          `db:"blk_nonce" json:"nonce"`
		Round           uint64          `db:"blk_round" json:"round"`
		Shard           uint64          `db:"blk_shard" json:"shard"`
		NumTxs          uint64          `db:"blk_num_txs" json:"num_txs"`
		Epoch           uint64          `db:"blk_epoch" json:"epoch"`
		Status          string          `db:"blk_status" json:"status"`
		PrevBlockHash   string          `db:"blk_prev_block_hash" json:"prev_block_hash"`
		AccumulatedFees decimal.Decimal `db:"blk_accumulated_fees" json:"accumulated_fees"`
		DeveloperFees   decimal.Decimal `db:"blk_developer_fees" json:"developer_fees"`
		CreatedAt       Time            `db:"blk_created_at" json:"created_at"`
	}
	MiniBlock struct {
		Hash              string `db:"mlk_hash" json:"hash"`
		ReceiverBlockHash string `db:"mlk_receiver_block_hash" json:"receiver_block_hash"` // can be empty
		ReceiverShard     uint64 `db:"mlk_receiver_shard" json:"receiver_shard"`
		SenderBlockHash   string `db:"mlk_sender_block_hash" json:"sender_block_hash"`
		SenderShard       uint64 `db:"mlk_sender_shard" json:"sender_shard"`
		Type              string `db:"mlk_type" json:"type"` // type == TxBlock
		CreatedAt         Time   `db:"mlk_created_at" json:"created_at"`
	}
)
