package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
)

func (db Postgres) CreateToken(token dmodels.Token) error {
	q := squirrel.Insert(dmodels.TokensTable).SetMap(map[string]interface{}{
		"tkn_identity":   token.Identity,
		"tkn_name":       token.Name,
		"tkn_type":       token.Type,
		"tkn_owner":      token.Owner,
		"tkn_supply":     token.Supply,
		"tkn_decimals":   token.Decimals,
		"tkn_properties": token.Properties,
		"tkn_roles":      token.Roles,
	})
	_, err := db.insert(q)
	return err
}

func (db *Postgres) UpdateToken(token dmodels.Token) error {
	q := squirrel.Update(dmodels.TokensTable).
		Where(squirrel.Eq{"tkn_identity": token.Identity}).
		SetMap(map[string]interface{}{
			"tkn_identity":   token.Identity,
			"tkn_name":       token.Name,
			"tkn_type":       token.Type,
			"tkn_owner":      token.Owner,
			"tkn_supply":     token.Supply,
			"tkn_decimals":   token.Decimals,
			"tkn_properties": token.Properties,
			"tkn_roles":      token.Roles,
		})
	return db.update(q)
}

func (db Postgres) GetTokens(filter filters.Tokens) (tokens []dmodels.Token, err error) {
	q := squirrel.Select("*").From(dmodels.TokensTable)
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset() != 0 {
		q = q.Offset(filter.Offset())
	}
	err = db.find(&tokens, q)
	return tokens, err
}

func (db Postgres) GetTokensCount(filter filters.Tokens) (total uint64, err error) {
	q := squirrel.Select("count(*) as total").From(dmodels.TokensTable)
	err = db.first(&total, q)
	return total, err
}

func (db Postgres) GetToken(ident string) (token dmodels.Token, err error) {
	q := squirrel.Select("*").From(dmodels.TokensTable).Where(squirrel.Eq{"tkn_identity": ident})
	err = db.first(&token, q)
	return token, err
}

func (db Postgres) CreateNFTCollection(collection dmodels.NFTCollection) error {
	q := squirrel.Insert(dmodels.NFTCollectionsTable).SetMap(map[string]interface{}{
		"nfc_name":       collection.Name,
		"nfc_identity":   collection.Identity,
		"nfc_owner":      collection.Owner,
		"nfc_type":       collection.Type,
		"nfc_properties": collection.Properties,
		"nfc_created_at": collection.CreatedAt,
	})
	_, err := db.insert(q)
	return err
}

func (db *Postgres) UpdateNFTCollection(collection dmodels.NFTCollection) error {
	q := squirrel.Update(dmodels.NFTCollectionsTable).
		Where(squirrel.Eq{"nfc_identity": collection.Identity}).
		SetMap(map[string]interface{}{
			"nfc_name":       collection.Name,
			"nfc_owner":      collection.Owner,
			"nfc_type":       collection.Type,
			"nfc_properties": collection.Properties,
			"nfc_created_at": collection.CreatedAt,
		})
	return db.update(q)
}

func (db Postgres) GetNFTCollections(filter filters.NFTCollections) (collections []dmodels.NFTCollection, err error) {
	q := squirrel.Select("*").From(dmodels.NFTCollectionsTable).OrderBy("nfc_created_at desc")
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset() != 0 {
		q = q.Offset(filter.Offset())
	}
	err = db.find(&collections, q)
	return collections, err
}

func (db Postgres) GetNFTCollectionsTotal(filter filters.NFTCollections) (total uint64, err error) {
	q := squirrel.Select("count(*) as total").From(dmodels.NFTCollectionsTable)
	err = db.first(&total, q)
	return total, err
}

func (db Postgres) GetNFTCollection(ident string) (collection dmodels.NFTCollection, err error) {
	q := squirrel.Select("*").From(dmodels.NFTCollectionsTable).Where(squirrel.Eq{"nfc_identity": ident})
	err = db.first(&collection, q)
	return collection, err
}
