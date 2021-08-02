package dailystats

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/shopspring/decimal"
)

func (ds *DailyStats) GetEconomics() (map[string]decimal.Decimal, error) {
	data, err := ds.node.GetNetworkEconomics()
	if err != nil {
		return nil, fmt.Errorf("node.GetNetworkEconomics: %s", err.Error())
	}
	return map[string]decimal.Decimal{
		TotalFeeKey:    node.ValueToEGLD(data.ErdTotalFees),
		TotalStakeKey:  node.ValueToEGLD(data.ErdTotalBaseStakedValue),
		TotalSupplyKey: node.ValueToEGLD(data.ErdTotalSupply),
		TopUpAmountKey: node.ValueToEGLD(data.ErdTotalTopUpValue),
	}, nil
}
