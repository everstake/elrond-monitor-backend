package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/services/node"
	"github.com/everstake/elrond-monitor-backend/smodels"
)

func (s *ServiceFacade) GetEpoch() (epoch smodels.Epoch, err error) {
	status, err := s.node.GetNetworkStatus(node.MetaChainShardIndex)
	if err != nil {
		return epoch, fmt.Errorf("node.GetNetworkStatus: %s", err.Error())
	}
	return smodels.Epoch{
		CurrentRound:   status.ErdCurrentRound,
		EpochNumber:    status.ErdEpochNumber,
		Nonce:          status.ErdNonce,
		RoundsPerEpoch: status.ErdRoundsPerEpoch,
	}, nil
}
