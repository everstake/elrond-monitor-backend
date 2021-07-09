package smodels

type EpochStats struct {
	Number float64 `json:"number,omitempty"`

	EndsAfter float64 `json:"ends_after,omitempty"`

	TotalStake float64 `json:"total_stake,omitempty"`

	TotalDelegators float64 `json:"total_delegators,omitempty"`

	PriceShift float64 `json:"price_shift,omitempty"`
}
