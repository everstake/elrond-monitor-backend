package filters

import "fmt"

const defaultPageLimit = 10

type Pagination struct {
	Limit    uint64 `schema:"limit"`
	Page     uint64 `schema:"page"`
	maxLimit uint64
}

func (p *Pagination) Validate() error {
	if p.maxLimit != 0 && p.Limit > p.maxLimit {
		return fmt.Errorf("overflow max limit")
	}
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = defaultPageLimit
	}
	return nil
}

func (p *Pagination) Offset() uint64 {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Limit * (p.Page - 1)
}

func (p *Pagination) SetMaxLimit(limit uint64) {
	p.maxLimit = limit
}
