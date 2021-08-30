package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
)

func (db Postgres) CreateDelegations(delegations []dmodels.Delegation) error {
	if len(delegations) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.DelegationsTable).Columns(
		"dlg_tx_hash",
		"dlg_delegator",
		"dlg_validator",
		"dlg_amount",
		"dlg_created_at",
	)
	for _, dlg := range delegations {
		if dlg.TxHash == "" {
			return fmt.Errorf("field TxHash is empty")
		}
		if dlg.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			dlg.TxHash,
			dlg.Delegator,
			dlg.Validator,
			dlg.Amount,
			dlg.CreatedAt,
		)
	}
	q = q.Suffix("ON CONFLICT (dlg_tx_hash) DO NOTHING")
	_, err := db.insert(q)
	return err
}