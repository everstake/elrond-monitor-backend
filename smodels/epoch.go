package smodels

type Epoch struct {
	CurrentRound   uint64 `json:"current_round"`
	EpochNumber    uint64 `json:"epoch_number"`
	Nonce          uint64 `json:"nonce"`
	RoundsPerEpoch uint64 `json:"rounds_per_epoch"`
}
