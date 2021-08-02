package filters

type DailyStats struct {
	Key   string `schema:"-"`
	Limit uint64 `schema:"limit"`
}
