package model

import (
	"github.com/lock-free/gopcp"
	"github.com/lock-free/obrero/napool"
	"github.com/lock-free/obrero/utils"
	"time"
)

var pcpClient = gopcp.PcpClient{}

// MapModel, provides simple apis to access DB model
type MapModel struct {
	DB      string
	key     string
	naPools *napool.NAPools
}

func GetMapModel(naPools *napool.NAPools, DB string, key string) *MapModel {
	return &MapModel{
		DB:      DB,
		key:     key,
		naPools: naPools,
	}
}

func (m MapModel) Get(modelPointer interface{}) error {
	d, err := m.naPools.CallProxy("db_obrero", pcpClient.Call("getByKey", m.DB, m.key, 120), 120*time.Second)
	if err != nil {
		return err
	}
	return utils.ParseArg(d, modelPointer)
}

func (m MapModel) Set(v interface{}) (interface{}, error) {
	return m.naPools.CallProxy("db_obrero", pcpClient.Call("setByKey", m.DB, m.key, v, 120), 120*time.Second)
}

func (m MapModel) Delete() (interface{}, error) {
	return m.naPools.CallProxy("db_obrero", pcpClient.Call("deleteByKey", m.DB, m.key, 120), 120*time.Second)
}
