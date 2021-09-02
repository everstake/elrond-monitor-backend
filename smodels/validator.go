package smodels

import "github.com/shopspring/decimal"

type Identity struct {
	Identity     string          `json:"identity"`
	Name         string          `json:"name"`
	Avatar       string          `json:"avatar"`
	Description  string          `json:"description"`
	Locked       decimal.Decimal `json:"locked"`
	Rank         uint64          `json:"rank"`
	Score        uint64          `json:"score"`
	Stake        decimal.Decimal `json:"stake"`
	StakePercent float64         `json:"stake_percent"`
	TopUp        decimal.Decimal `json:"top_up"`
	Validators   uint64          `json:"validators"`
	AVGUptime    float64         `json:"avg_uptime"`
	Providers    []string        `json:"providers"`
}

type StakingProvider struct {
	Provider         string                   `json:"provider"`
	ServiceFee       decimal.Decimal          `json:"service_fee"`
	DelegationCap    decimal.Decimal          `json:"delegation_cap"`
	APR              decimal.Decimal          `json:"apr"`
	NumUsers         uint64                   `json:"num_users"`
	CumulatedRewards decimal.Decimal          `json:"cumulated_rewards"`
	Identity         string                   `json:"identity"`
	Name             string                   `json:"name"`
	NumNodes         uint64                   `json:"num_nodes"`
	Stake            decimal.Decimal          `json:"stake"`
	TopUp            decimal.Decimal          `json:"top_up"`
	Locked           decimal.Decimal          `json:"locked"`
	Featured         bool                     `json:"featured"`
	Validator        StakingProviderValidator `json:"validator"`
}

type StakingProviderValidator struct {
	Name         string          `json:"name"`
	Locked       decimal.Decimal `json:"locked"`
	StakePercent float64         `json:"stake_percent"`
	Nodes        uint64          `json:"nodes"`
}

type SourceStakingProvider struct {
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

type IdentityKeybase struct {
	Status struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"status"`
	Them struct {
		ID     string `json:"id"`
		Basics struct {
			Username      string `json:"username"`
			Ctime         int    `json:"ctime"`
			Mtime         int    `json:"mtime"`
			IDVersion     int    `json:"id_version"`
			TrackVersion  int    `json:"track_version"`
			LastIDChange  int    `json:"last_id_change"`
			UsernameCased string `json:"username_cased"`
			Status        int    `json:"status"`
			Salt          string `json:"salt"`
			EldestSeqno   int    `json:"eldest_seqno"`
		} `json:"basics"`
		Profile struct {
			Mtime    interface{} `json:"mtime"`
			FullName string      `json:"full_name"`
			Location interface{} `json:"location"`
			Bio      string      `json:"bio"`
		} `json:"profile"`
		PublicKeys struct {
			AllBundles           []string `json:"all_bundles"`
			Subkeys              []string `json:"subkeys"`
			Sibkeys              []string `json:"sibkeys"`
			EldestKid            string   `json:"eldest_kid"`
			EldestKeyFingerprint string   `json:"eldest_key_fingerprint"`
		} `json:"public_keys"`
		ProofsSummary struct {
			ByPresentationGroup struct {
			} `json:"by_presentation_group"`
			BySigID struct {
			} `json:"by_sig_id"`
			HasWeb bool `json:"has_web"`
		} `json:"proofs_summary"`
		CryptocurrencyAddresses struct {
		} `json:"cryptocurrency_addresses"`
		Pictures struct {
			Primary struct {
				URL    string      `json:"url"`
				Source interface{} `json:"source"`
			} `json:"primary"`
		} `json:"pictures"`
		Sigs struct {
			Last struct {
				SigID       string `json:"sig_id"`
				Seqno       int    `json:"seqno"`
				PayloadHash string `json:"payload_hash"`
			} `json:"last"`
		} `json:"sigs"`
		Stellar struct {
			Hidden  bool `json:"hidden"`
			Primary struct {
			} `json:"primary"`
		} `json:"stellar"`
	} `json:"them"`
}
