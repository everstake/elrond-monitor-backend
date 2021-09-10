package filters

const (
	NodesSortByOnline = "online"
	NodesSortByShard  = "shard"
)

type Nodes struct {
	Pagination
	Identity string `schema:"identity"`
	Provider string `schema:"provider"`
	SortBy   string `schema:"sort_by"`
	Desc     bool   `schema:"desc"`
}
