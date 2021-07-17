package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
)

func (db Postgres) CreateRewards(rewards []dmodels.Reward) error {
	if len(rewards) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.RewardsTable).Columns(
		"rwd_tx_hash",
		"rwd_hyperblock_id",
		"rwd_receiver_address",
		"rwd_amount",
		"rwd_created_at",
	)
	for _, rwd := range rewards {
		if rwd.TxHash == "" {
			return fmt.Errorf("field TxHash is empty")
		}
		if rwd.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			rwd.TxHash,
			rwd.HypeblockID,
			rwd.ReceiverAddress,
			rwd.Amount,
			rwd.CreatedAt,
		)
	}
	q = q.Suffix("ON CONFLICT (rwd_tx_hash) DO NOTHING")
	_, err := db.insert(q)
	return err
}
