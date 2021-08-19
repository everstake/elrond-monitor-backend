package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
)

func (db Postgres) CreateDailyStats(stats []dmodels.DailyStat) error {
	if len(stats) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.DailyStatsTable).Columns(
		"das_title",
		"das_value",
		"das_created_at",
	)
	for _, s := range stats {
		if s.Title == "" {
			return fmt.Errorf("field Title is empty")
		}
		if s.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt is empty")
		}
		q = q.Values(
			s.Title,
			s.Value,
			s.CreatedAt,
		)
	}
	_, err := db.insert(q)
	return err
}

func (db Postgres) GetDailyStatsRange(filter filters.DailyStats) (items []dmodels.DailyStat, err error) {
	q := squirrel.Select("*").
		From(dmodels.DailyStatsTable).
		OrderBy("das_created_at desc").
		Where(squirrel.Eq{"das_title": filter.Key})
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if !filter.From.IsZero() {
		q = q.Where(squirrel.GtOrEq{"das_created_at": filter.From})
	}
	if !filter.To.IsZero() {
		q = q.Where(squirrel.LtOrEq{"das_created_at": filter.To})
	}
	err = db.find(&items, q)
	return items, err
}
