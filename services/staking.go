package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/smodels"
)

func (s *ServiceFacade) GetStakeEvents(filter filters.StakeEvents) (page smodels.Pagination, err error) {
	items, err := s.dao.GetStakeEvents(filter)
	if err != nil {
		return page, fmt.Errorf("dao.GetStakeEvents: %s", err.Error())
	}
	total, err := s.dao.GetStakeEventsTotal(filter)
	if err != nil {
		return page, fmt.Errorf("dao.GetStakeEventsTotal: %s", err.Error())
	}
	events := make([]smodels.StakeEvent, len(items))
	for i, item := range items {
		events[i] = smodels.StakeEvent{
			TxHash:    item.TxHash,
			Type:      item.Type,
			Validator: item.Validator,
			Delegator: item.Delegator,
			Epoch:     item.Epoch,
			Amount:    item.Amount,
			CreatedAt: smodels.NewTime(item.CreatedAt),
		}
	}
	return smodels.Pagination{
		Items: events,
		Count: total,
	}, nil
}
