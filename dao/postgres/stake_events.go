package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
)

func (db Postgres) CreateStakeEvents(events []dmodels.StakeEvent) error {
	if len(events) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.StakeEventsTable).Columns(
		"ste_tx_hash",
		"ste_type",
		"ste_validator",
		"ste_delegator",
		"ste_epoch",
		"ste_amount",
		"ste_created_at",
	)
	for _, e := range events {
		if e.TxHash == "" {
			return fmt.Errorf("field TxHash is empty")
		}
		if e.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			e.TxHash,
			e.Type,
			e.Validator,
			e.Delegator,
			e.Epoch,
			e.Amount,
			e.CreatedAt,
		)
	}
	q = q.Suffix("ON CONFLICT (ste_tx_hash) DO NOTHING")
	_, err := db.insert(q)
	return err
}
