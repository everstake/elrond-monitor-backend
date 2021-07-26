package filters

type Blocks struct {
	Shard []uint64 `schema:"shard"`
	Nonce uint64   `schema:"nonce"`
	Pagination
}

type MiniBlocks struct {
	ParentBlockHash string
}
