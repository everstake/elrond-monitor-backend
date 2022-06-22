package filters

type Tokens struct {
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
	Pagination
}
