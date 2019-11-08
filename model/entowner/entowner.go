package entowner

import (
	"errors"
	"github.com/lock-free/gopcp"
	"github.com/lock-free/obrero/model"
	"github.com/lock-free/obrero/napool"
	"github.com/satori/go.uuid"
	"time"
)

var pcpClient = gopcp.PcpClient{}

// define CRUD for owner level entity

type EntOnwer struct {
	RelKey    string
	EntityKey string
}

func (eo EntOnwer) checkPermission(naPools *napool.NAPools, oid string, eid string) error {
	return model.CheckRel(naPools, eo.RelKey, oid, eid)
}

func (eo EntOnwer) SetEnt(naPools *napool.NAPools, oid, eid string, entity map[string]interface{}) (interface{}, error) {
	var err error

	if eid == "" { // create new entity
		eid = uuid.NewV4().String()
		_, err = naPools.CallProxy("model_obrero", pcpClient.Call("setEntity", eo.EntityKey, eid, entity), 120*time.Second)
		if err != nil {
			return nil, err
		}
		// set a new relationship for owner and entity
		return naPools.CallProxy("model_obrero", pcpClient.Call("setRel", eo.RelKey, oid, eid), 120*time.Second)
	} else {
		// check permission
		if err = eo.checkPermission(naPools, oid, eid); err != nil {
			return nil, err
		}
		// update
		return naPools.CallProxy("model_obrero", pcpClient.Call("setEntity", eo.EntityKey, eid, entity), 120*time.Second)
	}
}

func (eo EntOnwer) DeleteEnt(naPools *napool.NAPools, oid, eid string) (interface{}, error) {
	if err := eo.checkPermission(naPools, oid, eid); err != nil {
		return nil, err
	}
	// delete relationship
	_, err := naPools.CallProxy("model_obrero", pcpClient.Call("deleteRel", eo.EntityKey, oid, eid), 120*time.Second)
	if err != nil {
		return nil, err
	}
	// delete enti
	return naPools.CallProxy("model_obrero", pcpClient.Call("deleteEntity", eo.EntityKey, eid), 120*time.Second)
}

func (eo EntOnwer) GetEnt(naPools *napool.NAPools, oid, eid string) (map[string]interface{}, error) {
	if err := eo.checkPermission(naPools, oid, eid); err != nil {
		return nil, err
	}
	v, err := naPools.CallProxy("model_obrero", pcpClient.Call("getEntity", eo.EntityKey, eid), 120*time.Second)
	if err != nil {
		return nil, err
	}

	ent, ok := v.(map[string]interface{})

	if !ok {
		return nil, errors.New("type error for entity, expect map[string]interface{}")
	}

	return ent, nil
}

func (eo EntOnwer) GetEnts(naPools *napool.NAPools, oid string) (interface{}, error) {
	return naPools.CallProxy("model_obrero", pcpClient.Call("getRels", eo.RelKey, oid), 120*time.Second)
}
