package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
)

func (db *Postgres) UpdateStorageValue(item dmodels.StorageItem) error {
	q := squirrel.Update(dmodels.StorageTable).
		Where(squirrel.Eq{"stg_key": item.Key}).
		SetMap(map[string]interface{}{
			"stg_value": item.Value,
		})
	return db.update(q)
}

func (db *Postgres) GetStorageValue(key string) (value string, err error) {
	q := squirrel.Select("stg_value").
		From(dmodels.StorageTable).
		Where(squirrel.Eq{"stg_key": key})
	err = db.first(&value, q)
	return value, err
}
