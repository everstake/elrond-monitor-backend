package filters

type Transactions struct {
	Pagination
	Address   string `schema:"address"`
	MiniBlock string `schema:"mini_block"`
}

type Operations struct {
	Token string `schema:"token"`
	Pagination
}
