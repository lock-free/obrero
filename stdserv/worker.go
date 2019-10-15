package stdserv

import (
	"errors"
	"github.com/lock-free/gopcp"
	"github.com/lock-free/gopcp_stream"
	"github.com/lock-free/obrero"
	"github.com/lock-free/obrero/mids"
	"github.com/lock-free/obrero/napool"
	"github.com/lock-free/obrero/utils"
	"time"
)

// define standard worker
// (1) read configuration from /data/app.json
// (2) default middlewares

const (
	DEFAULT_APP_CONFIG = "/data/app.json"
)

var DEFAULT_WORKER_START_CONF = obrero.WorkerStartConf{
	PoolSize:            2,
	Duration:            20 * time.Second,
	RetryDuration:       20 * time.Second,
	NAGetClientMaxRetry: 3,
}

type StdWorkerConfig struct {
	// TODO read ServiceName from env
	ServiceName       string
	AppConfigFilePath *string
	WorkerStartConf   *obrero.WorkerStartConf
}

// (pool, config pointer, stream)
type GetBoxFuncMap = func(*napool.NAPools, *gopcp_stream.StreamServer) map[string]*gopcp.BoxFunc

// appConfig: pointer of appConfig
func StartStdWorker(appConfig interface{}, getBoxFuncMap GetBoxFuncMap, stdWorkerConfig StdWorkerConfig) {
	// read config from config file
	appConfigFilePath := DEFAULT_APP_CONFIG
	if stdWorkerConfig.AppConfigFilePath != nil {
		appConfigFilePath = *stdWorkerConfig.AppConfigFilePath
	}

	err := utils.ReadJson(appConfigFilePath, appConfig)
	if err != nil {
		panic(err)
	}

	// start worker
	var naPools napool.NAPools

	var workerStartConf = DEFAULT_WORKER_START_CONF
	if stdWorkerConfig.WorkerStartConf != nil {
		workerStartConf = *stdWorkerConfig.WorkerStartConf
	}

	naPools = obrero.StartWorker(func(s *gopcp_stream.StreamServer) *gopcp.Sandbox {
		boxFuncMap := getBoxFuncMap(&naPools, s)

		for key, boxFunc := range boxFuncMap {
			// log function
			boxFunc.Fun = mids.LogMid(key, boxFunc.Fun)
		}

		// default get service type function
		if _, ok := boxFuncMap["getServiceType"]; !ok {
			if stdWorkerConfig.ServiceName == "" {
				panic(errors.New("missing service name in std service"))
			}
			boxFuncMap["getServiceType"] = gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
				return stdWorkerConfig.ServiceName, nil
			})
		}

		return gopcp.GetSandbox(boxFuncMap)
	}, workerStartConf)

	utils.RunForever()
}
