package filters

const (
	NodesSortByStatus = "online"
	NodesSortByShard  = "shard"
)

type Nodes struct {
	Pagination
	Identity string   `schema:"identity"`
	Provider string   `schema:"provider"`
	Status   []uint64 `schema:"status"`
	Shard    []uint64 `schema:"shard"`
}
