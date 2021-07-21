package services

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/elrond-monitor-backend/dao/dmodels"
)

func (s *ServiceFacade) setCache(key string, item interface{}) error {
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("json.Marshal: %s", err.Error())
	}
	err = s.dao.UpdateStorageValue(dmodels.StorageItem{
		Key:   key,
		Value: string(data),
	})
	if err != nil {
		return fmt.Errorf("dao.UpdateStorageValue: %s", err.Error())
	}
	return nil
}

func (s *ServiceFacade) getCache(key string, dst interface{}) error {
	value, err := s.dao.GetStorageValue(key)
	if err != nil {
		return fmt.Errorf("dao.GetStorageValue: %s", err.Error())
	}
	err = json.Unmarshal([]byte(value), dst)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	return nil
}
