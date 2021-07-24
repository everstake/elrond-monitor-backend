package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
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

