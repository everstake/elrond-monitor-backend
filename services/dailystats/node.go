package dailystats

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func (ds *DailyStats) GetEconomics() (map[string]decimal.Decimal, error) {
	data, err := ds.node.GetNetworkEconomics()
	if err != nil {
		return nil, fmt.Errorf("node.GetNetworkEconomics: %s", err.Error())
	}
	return map[string]decimal.Decimal{
		TotalFeeKey:    data.ErdTotalFees,
		TotalStakeKey:  data.ErdTotalBaseStakedValue,
		TotalSupplyKey: data.ErdTotalSupply,
		TopUpAmountKey: data.ErdTotalTopUpValue,
	}, nil
}
