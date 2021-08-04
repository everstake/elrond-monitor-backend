package smodels

type Epoch struct {
	CurrentRound   uint64  `json:"current_round"`
	EpochNumber    uint64  `json:"epoch_number"`
	Nonce          uint64  `json:"nonce"`
	RoundsPerEpoch uint64  `json:"rounds_per_epoch"`
	Percent        float64 `json:"percent"`
	Left           uint64  `json:"left"`
	Start          Time    `json:"start"`
}
