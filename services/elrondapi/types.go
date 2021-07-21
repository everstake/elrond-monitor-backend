package elrondapi

import "github.com/shopspring/decimal"

type (
	Stats struct {
		Shards         uint64 `json:"shards"`
		Blocks         uint64 `json:"blocks"`
		Accounts       uint64 `json:"accounts"`
		Transactions   uint64 `json:"transactions"`
		RefreshRates   uint64 `json:"refreshRate"`
		Epoch          uint64 `json:"epoch"`
		RoundsPerEpoch uint64 `json:"roundsPerEpoch"`
	}

	Identity struct {
		Identity     string                 `json:"identity"`
		Name         string                 `json:"name"`
		Description  string                 `json:"description"`
		Avatar       string                 `json:"avatar"`
		Score        uint64                 `json:"score"`
		Validators   uint64                 `json:"validators"`
		Stake        string                 `json:"stake"`
		TopUp        string                 `json:"topUp"`
		Locked       string                 `json:"locked"`
		Distribution map[string]interface{} `json:"distribution"`
		Providers    []string               `json:"providers"`
		StakePercent decimal.Decimal        `json:"stakePercent"`
		Rank         uint64                 `json:"rank"`
	}

	Economics struct {
		TotalSupply       uint64 `json:"totalSupply"`
		CirculatingSupply uint64 `json:"circulatingSupply"`
		Staked            uint64 `json:"staked"`
	}

	AccountDelegation struct {
		UserWithdrawOnlyStake    decimal.Decimal `json:"userWithdrawOnlyStake"`
		UserWaitingStake         decimal.Decimal `json:"userWaitingStake"`
		UserActiveStake          decimal.Decimal `json:"userActiveStake"`
		UserUnstakedStake        decimal.Decimal `json:"userUnstakedStake"`
		UserDeferredPaymentStake decimal.Decimal `json:"userDeferredPaymentStake"`
		ClaimableRewards         decimal.Decimal `json:"claimableRewards"`
	}
)
