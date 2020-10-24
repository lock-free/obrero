package mysql

import (
	"errors"
	"time"

	"github.com/LoveKino/model_sql"
	"github.com/lock-free/gopcp"
	"github.com/lock-free/obrero/napool"
	"github.com/lock-free/obrero/utils"
)

var (
	EMPTY_ERROR = errors.New("empty error")
)

var pcpClient = gopcp.PcpClient{}

// TODO options
func Upsert(naPools *napool.NAPools, dbName, tableName string, model interface{}, primary []string) (interface{}, error) {
	// get upsert sql
	sql, err := model_sql.GetUpsertSQL("payment", model, primary)
	if err != nil {
		return nil, err
	}

	// execute sql
	return naPools.CallProxy("data_service", pcpClient.Call("execMysql", dbName, sql, 120), 120*time.Second)
}

func GetModelsByFields(naPools *napool.NAPools, dbName, tableName string, fields map[string]interface{}) ([]map[string]interface{}, error) {
	sql, err := model_sql.GetByFieldsSQL(tableName, fields)

	if err != nil {
		return nil, err
	}

	// query sql
	v, err := naPools.CallProxy("data_service", pcpClient.Call("queryMysql", dbName, sql, 120), 120*time.Second)
	if err != nil {
		return nil, err
	}

	//
	var models []map[string]interface{}
	err = utils.ParseArg(v, &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func GetFirstByFields(naPools *napool.NAPools, dbName, tableName string, fields map[string]interface{}, ptr interface{}) error {
	sql, err := model_sql.GetFirstByFieldsSQL(tableName, fields)

	if err != nil {
		return err
	}

	// query sql
	v, err := naPools.CallProxy("data_service", pcpClient.Call("queryMysql", dbName, sql, 120), 120*time.Second)
	if err != nil {
		return err
	}

	//
	var models []map[string]interface{}
	err = utils.ParseArg(v, &models)
	if err != nil {
		return err
	}

	if len(models) == 0 {
		return EMPTY_ERROR
	}

	return utils.ParseArg(models[0], ptr)
}

func CountByFields(naPools *napool.NAPools, dbName, tableName string, fields map[string]interface{}) (uint64, error) {
	sql, err := model_sql.CountByFieldsSQL(tableName, fields)
	if err != nil {
		return 0, err
	}

	// query sql
	v, err := naPools.CallProxy("data_service", pcpClient.Call("queryMysql", dbName, sql, 120), 120*time.Second)
	if err != nil {
		return 0, err
	}

	// parse models
	var models []map[string]interface{}
	err = utils.ParseArg(v, &models)
	if err != nil {
		return 0, err
	}

	if len(models) == 0 {
		return 0, EMPTY_ERROR
	}

	// get count
	for _, v := range models[0] {
		c, ok := v.(uint64)
		if !ok {
			return 0, errors.New("unexpect value type for count(*)")
		}

		return c, nil
	}

	return 0, nil
}

func ExistsByFields(naPools *napool.NAPools, dbName, tableName string, fields map[string]interface{}) (bool, error) {
	c, err := CountByFields(naPools, dbName, tableName, fields)
	if err != nil {
		return false, err
	}

	if c == 0 {
		return false, nil
	}
	return true, nil
}
