package swagger

type Validator struct {
	Name string `json:"name,omitempty"`

	Stake float64 `json:"stake,omitempty"`

	CumulativeStake float64 `json:"cumulative_stake,omitempty"`

	Nodes float64 `json:"nodes,omitempty"`
}
