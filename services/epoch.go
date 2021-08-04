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
	percent := float64(status.ErdNoncesPassedInCurrentEpoch) / float64(s.networkConfig.ErdRoundsPerEpoch) * 100
	left := (s.networkConfig.ErdRoundsPerEpoch - status.ErdNoncesPassedInCurrentEpoch) * s.networkConfig.ErdRoundDuration
	start := time.Now().Add(- time.Duration(s.networkConfig.ErdRoundDuration*status.ErdNoncesPassedInCurrentEpoch) * time.Millisecond)
	return smodels.Epoch{
		CurrentRound:   status.ErdCurrentRound,
		EpochNumber:    status.ErdEpochNumber,
		Nonce:          status.ErdNonce,
		RoundsPerEpoch: status.ErdRoundsPerEpoch,
		Percent:        percent,
		Left:           left / 1000,
		Start:          smodels.NewTime(start),
	}, nil
}
