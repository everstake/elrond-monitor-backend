package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/smodels"
)

func (s *ServiceFacade) GetDailyStats(filter filters.DailyStats) (items []smodels.RangeItem, err error) {
	dItems, err := s.dao.GetDailyStatsRange(filter)
	if err != nil {
		return items, fmt.Errorf("dao.GetDailyStatsRange: %s", err.Error())
	}
	items = make([]smodels.RangeItem, len(dItems))
	for i, it := range dItems {
		items[i] = smodels.RangeItem{
			Value: it.Value,
			Time:  smodels.NewTime(it.CreatedAt),
		}
	}
	return items, nil
}
