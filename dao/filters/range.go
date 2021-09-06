package filters

import "github.com/everstake/elrond-monitor-backend/smodels"

type DailyStats struct {
	Key   string       `schema:"-"`
	Limit uint64       `schema:"limit"`
	From  smodels.Time `schema:"from"`
	To    smodels.Time `schema:"to"`
}
