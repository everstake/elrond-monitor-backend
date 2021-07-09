package smodels

type Account struct {
	Address string `json:"address,omitempty"`

	Balance float64 `json:"balance,omitempty"`

	Delegated float64 `json:"delegated,omitempty"`

	Undelegated float64 `json:"undelegated,omitempty"`

	RewardsClaimed float64 `json:"rewards_claimed,omitempty"`

	StakingProvider float64 `json:"staking_provider,omitempty"`
}
