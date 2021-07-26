package smodels

import "github.com/shopspring/decimal"

type Account struct {
	Address          string          `json:"address"`
	Balance          decimal.Decimal `json:"balance"`
	Nonce            uint64          `json:"nonce"`
	Delegated        decimal.Decimal `json:"delegated"`
	Undelegated      decimal.Decimal `json:"undelegated"`
	RewardsClaimed   decimal.Decimal `json:"rewards_claimed"`
	ClaimableRewards decimal.Decimal `json:"claimable_rewards"`
	StakingProvider  []string        `json:"staking_provider"`
	CreatedAt        Time            `json:"created_at"`
}
