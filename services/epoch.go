package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
	"time"
)

func (s *ServiceFacade) GetEpoch() (epoch smodels.Epoch, err error) {
	status, err := s.node.GetNetworkStatus(node.MetaChainShardIndex)
	if err != nil {
		return epoch, fmt.Errorf("node.GetNetworkStatus: %s", err.Error())
	}
	if s.networkConfig.ErdRoundsPerEpoch == 0 {
		return epoch, fmt.Errorf("RoundsPerEpoch is zero")
	}
	roundsLeft := status.ErdCurrentRound - s.networkConfig.ErdRoundsPerEpoch*status.ErdEpochNumber
	percent := (float64(roundsLeft) / float64(s.networkConfig.ErdRoundsPerEpoch)) * 100
	start := s.networkConfig.ErdStartTime*1000 + status.ErdEpochNumber*s.networkConfig.ErdRoundsPerEpoch*s.networkConfig.ErdRoundDuration
	finalRound := status.ErdCurrentRound + roundsLeft
	end := start + s.networkConfig.ErdRoundDuration*finalRound
	return smodels.Epoch{
		CurrentRound:   status.ErdCurrentRound,
		EpochNumber:    status.ErdEpochNumber,
		Nonce:          status.ErdNonce,
		RoundsPerEpoch: status.ErdRoundsPerEpoch,
		Percent:        percent,
		Start:          smodels.NewTime(time.Unix(int64(start/1000), 0)),
		End:            smodels.NewTime(time.Unix(int64(end/1000), 0)),
	}, nil
}
