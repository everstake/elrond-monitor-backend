package parser

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
	"github.com/shopspring/decimal"
)

func (p *Parser) loadStates() error {
	delegations, err := p.dao.GetDelegationState()
	if err != nil {
		return fmt.Errorf("dao.GetDelegationState: %s", err.Error())
	}
	for _, d := range delegations {
		if _, ok := p.delegations[d.Delegator]; !ok {
			p.delegations[d.Delegator] = make(map[string]decimal.Decimal)
		}
		p.delegations[d.Delegator][d.Validator] = d.Amount
	}
	return nil
}

func (p *Parser) updateStakeStates(events []dmodels.StakeEvent) {
	p.mu.Lock()
	for _, event := range events {
		if event.Type == dmodels.DelegateStakeEventType || event.Type == dmodels.UnDelegateStakeEventType {
			_, ok := p.delegations[event.Delegator]
			if !ok {
				p.delegations[event.Delegator] = make(map[string]decimal.Decimal)
			}
			v := p.delegations[event.Delegator][event.Validator]
			amount := v.Add(event.Amount)
			p.delegations[event.Delegator][event.Validator] = amount
			if amount.IsZero() {
				delete(p.delegations[event.Delegator], event.Validator)
			}
		}
	}
	p.mu.Unlock()
}

func (p *Parser) GetDelegations(delegator string) map[string]decimal.Decimal {
	res := make(map[string]decimal.Decimal)
	p.mu.RLock()
	m := p.delegations[delegator]
	for validator, amount := range m {
		res[validator] = amount
	}
	p.mu.RUnlock()
	return res
}
