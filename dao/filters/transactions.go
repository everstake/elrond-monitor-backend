package filters

type Transactions struct {
	Pagination
	Address string `schema:"address"`
}
