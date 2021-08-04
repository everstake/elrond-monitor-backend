package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
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

func (db Postgres) GetDelegationState() (items []dmodels.StakeState, err error) {
	q := squirrel.Select("ste_validator as validator", "ste_delegator as delegator", "sum(ste_amount) as amount").
		From(dmodels.StakeEventsTable).
		Where(squirrel.Eq{"ste_type": []string{dmodels.DelegateStakeEventType, dmodels.UnDelegateStakeEventType}}).
		GroupBy("ste_delegator", "ste_validator").
		Having("sum(ste_amount) > 0")
	err = db.find(&items, q)
	return items, err
}

func (db Postgres) GetStakeState() (items []dmodels.StakeState, err error) {
	q := squirrel.Select("ste_validator as validator", "ste_delegator as delegator", "sum(ste_amount) as amount").
		From(dmodels.StakeEventsTable).
		Where(squirrel.Eq{"ste_type": []string{dmodels.StakeStakeEventType, dmodels.UnStakeEventType}}).
		GroupBy("ste_delegator", "ste_validator").
		Having("sum(ste_amount) > 0")
	err = db.find(&items, q)
	return items, err
}

func (db Postgres) GetStakeEvents(filter filters.StakeEvents) (items []dmodels.StakeEvent, err error) {
	q := squirrel.Select("*").
		From(dmodels.StakeEventsTable).
		OrderBy("ste_created_at desc")
	if len(filter.Delegator) > 0 {
		q = q.Where(squirrel.Eq{"ste_delegator": filter.Delegator})
	}
	if len(filter.Validator) > 0 {
		q = q.Where(squirrel.Eq{"ste_validator": filter.Validator})
	}
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset() != 0 {
		q = q.Offset(filter.Offset())
	}
	err = db.find(&items, q)
	return items, err
}

func (db Postgres) GetStakeEventsTotal(filter filters.StakeEvents) (total uint64, err error) {
	q := squirrel.Select("count(*)").
		From(dmodels.StakeEventsTable)
	if len(filter.Delegator) > 0 {
		q = q.Where(squirrel.Eq{"ste_delegator": filter.Delegator})
	}
	if len(filter.Validator) > 0 {
		q = q.Where(squirrel.Eq{"ste_validator": filter.Validator})
	}
	err = db.first(&total, q)
	return total, err
}
