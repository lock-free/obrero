package mysql

import (
	"time"

	"github.com/LoveKino/model_sql"
	"github.com/lock-free/gopcp"
	"github.com/lock-free/obrero/napool"
)

var pcpClient = gopcp.PcpClient{}

// TODO options
func Upsert(naPools *napool.NAPools, dbName, tableName string, model interface{}, primary []string) (interface{}, error) {
	// get upsert sql
	sql, err := model_sql.GetUpsertSQL("payment", model, []string{"pay_id"})
	if err != nil {
		return nil, err
	}

	// execute sql
	return naPools.CallProxy("data_service", pcpClient.Call("execMysql", dbName, sql, 120), 120*time.Second)
}
