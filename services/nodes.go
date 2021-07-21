package services

import (
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/filters"
	"github.com/everstake/elrond-monitor-backend/services/node"
)

const nodesStorageKey = "nodes"

func (s *ServiceFacade) UpdateNodes() error {
	nodesStatus, err := s.node.GetHeartbeatStatus()
	if err != nil {
		return fmt.Errorf("node.GetHeartbeatStatus: %s", err.Error())
	}
	err = s.setCache(nodesStorageKey, nodesStatus)
	if err != nil {
		return fmt.Errorf("setCache: %s", err.Error())
	}
	return nil
}

func (s *ServiceFacade) GetNodes(filter filters.Nodes) (nodes []node.HeartbeatStatus, err error) {
	err = s.getCache(nodesStorageKey, &nodes)
	if err != nil {
		return nil, fmt.Errorf("getCache: %s", err.Error())
	}
	nodesLen := uint64(len(nodes))
	if filter.Limit*(filter.Page-1) > nodesLen {
		return nil, nil
	}
	maxIndex := filter.Page*filter.Limit - 1
	if nodesLen-1 < maxIndex {
		maxIndex = nodesLen - 1
	}
	return nodes, nil
}
