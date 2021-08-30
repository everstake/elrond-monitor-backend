package dmodels

const (
	StorageTable = "storage"

	NodesStorageKey            = "nodes"
	StakingProvidersStorageKey = "staking_providers"
	ValidatorsStorageKey       = "validators"
	StatsStorageKey            = "stats"
	ValidatorStatsStorageKey   = "validator_stats"
	ValidatorsMapStorageKey    = "validators_map"
	RankingStorageKey          = "ranking"
)

type StorageItem struct {
	Key   string `db:"stg_key"`
	Value string `db:"stg_value"`
}
