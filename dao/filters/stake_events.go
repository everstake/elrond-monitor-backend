package filters

type StakeEvents struct {
	Validator []string `schema:"validator"`
	Delegator []string `schema:"delegator"`
	Pagination
}
