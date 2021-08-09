package smodels

import "github.com/shopspring/decimal"

type Validator struct {
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
	StakePercent decimal.Decimal        `json:"stake_percent"`
	Rank         uint64                 `json:"rank"`
}

type StakingProvider struct {
	Identity struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Avatar      string `json:"avatar"`
		Description string `json:"description"`
		Location    string `json:"location"`
	} `json:"identity"`
	Contract                           string          `json:"contract"`
	ExplorerURL                        string          `json:"explorerURL"`
	Featured                           bool            `json:"featured"`
	Owner                              string          `json:"owner"`
	ServiceFee                         decimal.Decimal `json:"serviceFee"`
	MaxDelegationCap                   decimal.Decimal `json:"maxDelegationCap"`
	InitialOwnerFunds                  decimal.Decimal `json:"initialOwnerFunds"`
	AutomaticActivation                bool            `json:"automaticActivation"`
	WithDelegationCap                  bool            `json:"withDelegationCap"`
	ChangeableServiceFee               bool            `json:"changeableServiceFee"`
	CheckCapOnRedelegate               bool            `json:"checkCapOnRedelegate"`
	CreatedNonce                       uint64          `json:"createdNonce"`
	UnBondPeriod                       uint64          `json:"unBondPeriod"`
	Apr                                decimal.Decimal `json:"apr"`
	AprValue                           decimal.Decimal `json:"aprValue"`
	TotalActiveStake                   decimal.Decimal `json:"totalActiveStake"`
	TotalUnStaked                      decimal.Decimal `json:"totalUnStaked"`
	NumUsers                           uint64          `json:"numUsers"`
	NumNodes                           uint64          `json:"numNodes"`
	MaxDelegateAmountAllowed           decimal.Decimal `json:"maxDelegateAmountAllowed"`
	MaxRedelegateAmountAllowed         decimal.Decimal `json:"maxRedelegateAmountAllowed"`
	OwnerBelowRequiredBalanceThreshold bool            `json:"ownerBelowRequiredBalanceThreshold"`
}
