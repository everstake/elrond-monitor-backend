package filters

type Nodes struct {
	Pagination
	Identity string `schema:"identity"`
	Provider string `schema:"provider"`
}
