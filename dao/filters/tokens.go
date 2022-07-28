package filters

type Tokens struct {
	Identifier []string `schema:"identifier"`
	Pagination
}

type NFTTokens struct {
	Collection string `schema:"collection"`
	Pagination
}

type NFTCollections struct {
	Pagination
}

type ESDT struct {
	TokenIdentifier string `schema:"token_identifier"`
	Address         string `schema:"address"`
	Pagination
}
