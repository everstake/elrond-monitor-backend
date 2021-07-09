package filters

const defaultPageLimit = 100

type Pagination struct {
	Limit uint64 `schema:"limit"`
	Page  uint64 `schema:"page"`
}

func (p *Pagination) Validate() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = defaultPageLimit
	}
}

func (p *Pagination) Offset() uint64 {
	return p.Limit * (p.Page - 1)
}
