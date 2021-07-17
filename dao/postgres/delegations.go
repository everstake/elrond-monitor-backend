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

func (db Postgres) CreateStakes(stakes []dmodels.Stake) error {
	if len(stakes) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.StakesTable).Columns(
		"stk_tx_hash",
		"stk_validator",
		"stk_amount",
		"stk_created_at",
	)
	for _, s := range stakes {
		if s.TxHash == "" {
			return fmt.Errorf("field TxHash is empty")
		}
		if s.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			s.TxHash,
			s.Validator,
			s.Amount,
			s.CreatedAt,
		)
	}
	q = q.Suffix("ON CONFLICT (stk_tx_hash) DO NOTHING")
	_, err := db.insert(q)
	return err
}
