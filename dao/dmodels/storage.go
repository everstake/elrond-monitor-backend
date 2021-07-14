package dmodels

const StorageTable = "storage"

type StorageItem struct {
	Key   string `db:"stg_key"`
	Value string `db:"stg_value"`
}
