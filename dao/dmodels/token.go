package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	TokensTable         = "tokens"
	NFTCollectionsTable = "nft_collections"
)

const (
	FungibleESDT     = "FungibleESDT"
	NonFungibleESDT  = "NonFungibleESDT"
	SemiFungibleESDT = "SemiFungibleESDT"
	MetaESDT         = "MetaESDT"
)

type Token struct {
	Identity   string          `db:"tkn_identity"`
	Name       string          `db:"tkn_name"`
	Type       string          `db:"tkn_type"`
	Owner      string          `db:"tkn_owner"`
	Supply     decimal.Decimal `db:"tkn_supply"`
	Decimals   uint64          `db:"tkn_decimals"`
	Properties []byte          `db:"tkn_properties"`
	Roles      []byte          `db:"tkn_roles"`
	Operations uint64          `db:"tkn_operations"`
}

type NFTCollection struct {
	Name       string    `db:"nfc_name"`
	Identity   string    `db:"nfc_identity"`
	Owner      string    `db:"nfc_owner"`
	Type       string    `db:"nfc_type"`
	Properties []byte    `db:"nfc_properties"`
	CreatedAt  time.Time `db:"nfc_created_at"`
}
