package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const DailyStatsTable = "daily_stats"

type DailyStat struct {
	Title     string          `db:"das_title" json:"title"`
	Value     decimal.Decimal `db:"das_value" json:"value"`
	CreatedAt time.Time       `db:"das_created_at" json:"created_at"`
}
