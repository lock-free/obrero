package common

import (
	"github.com/lock-free/obrero/dt"
	"github.com/lock-free/obrero/napool"
)

// common business logic

// return entity if exists, otherwise, return nil
func GetEntityIfExist(naPools *napool.NAPools, entityKey string, entityID string) (interface{}, error) {
	// check if user is VIP
	v, err := naPools.SimpleCall("model_obrero", "hasEntity", entityKey, entityID)
	if err != nil {
		return nil, err
	}
	if dt.Falsy(v) {
		return nil, nil
	}
	return naPools.SimpleCall("model_obrero", "getEntity", entityKey, entityID)
}
