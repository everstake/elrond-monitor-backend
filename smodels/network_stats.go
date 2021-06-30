package swagger

type NetworkStats struct {
	TotalStake float64 `json:"total_stake,omitempty"`

	TotalDelegators float64 `json:"total_delegators,omitempty"`

	TotalAccounts float64 `json:"total_accounts,omitempty"`

	TopUpStake float64 `json:"top_up_stake,omitempty"`

	ActiveValidators float64 `json:"active_validators,omitempty"`

	StakingApr float64 `json:"staking_apr,omitempty"`

	BlockTime float64 `json:"block_time,omitempty"`
}
