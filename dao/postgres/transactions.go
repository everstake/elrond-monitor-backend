package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
)

func (db Postgres) CreateTransactions(transactions []dmodels.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.TransactionsTable).Columns(
		"trn_hash",
		"trn_status",
		"mlk_mini_block_hash",
		"trn_value",
		"trn_fee",
		"trn_sender",
		"trn_sender_shard",
		"trn_receiver",
		"trn_receiver_shard",
		"trn_gas_price",
		"trn_gas_used",
		"trn_nonce",
		"trn_data",
		"trn_created_at",
	)
	for _, tx := range transactions {
		if tx.Hash == "" {
			return fmt.Errorf("field Hash is empty")
		}
		if tx.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			tx.Hash,
			tx.Status,
			tx.MiniBlockHash,
			tx.Value,
			tx.Fee,
			tx.Sender,
			tx.SenderShard,
			tx.Receiver,
			tx.ReceiverShard,
			tx.GasPrice,
			tx.GasUsed,
			tx.Nonce,
			tx.Data,
			tx.CreatedAt,
		)
	}
	q = q.Suffix("ON CONFLICT (trn_hash) DO NOTHING")
	_, err := db.insert(q)
	return err
}

func (db Postgres) CreateSCResults(results []dmodels.SCResult) error {
	if len(results) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.SCResultsTable).Columns(
		"scr_hash",
		"trn_hash",
		"scr_from",
		"scr_to",
		"scr_value",
		"scr_data",
	)
	for _, r := range results {
		if r.TxHash == "" {
			return fmt.Errorf("field TxHash is empty")
		}
		q = q.Values(
			r.Hash,
			r.TxHash,
			r.To,
			r.From,
			r.Value,
			r.Data,
		)
	}
	q = q.Suffix("ON CONFLICT (scr_hash) DO NOTHING")
	_, err := db.insert(q)
	return err
}
