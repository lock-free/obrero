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
	SERVICE_STATE_FILE = "/data/__service_state.json"
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
	BeforeStartWorker func()
	Nonblocking       bool
}

// (pool, config pointer, stream)
type GetBoxFuncMap = func(*napool.NAPools, *WorkerState, *gopcp_stream.StreamServer) map[string]*gopcp.BoxFunc

// appConfig: pointer of appConfig
// appState: pointer of appState
func StartStdWorker(appConfig interface{}, appState interface{}, getBoxFuncMap GetBoxFuncMap, stdWorkerConfig StdWorkerConfig) napool.NAPools {
	// read config from config file
	appConfigFilePath := DEFAULT_APP_CONFIG
	if stdWorkerConfig.AppConfigFilePath != nil {
		appConfigFilePath = *stdWorkerConfig.AppConfigFilePath
	}

	err := utils.ReadJson(appConfigFilePath, appConfig)
	if err != nil {
		panic(err)
	}

	// read state from state file
	workerState, err := GetWorkerState(SERVICE_STATE_FILE)
	if err != nil {
		panic(err)
	}
	if appState != nil {
		err := utils.ParseArg(workerState.State.Data, appState)
		if err != nil {
			panic(err)
		}
		workerState.State.Data = appState
	}

	// start worker
	var naPools napool.NAPools

	var workerStartConf = DEFAULT_WORKER_START_CONF
	if stdWorkerConfig.WorkerStartConf != nil {
		workerStartConf = *stdWorkerConfig.WorkerStartConf
	}

	// before start worker
	if stdWorkerConfig.BeforeStartWorker != nil {
		stdWorkerConfig.BeforeStartWorker()
	}

	naPools = obrero.StartWorker(func(s *gopcp_stream.StreamServer) *gopcp.Sandbox {
		boxFuncMap := getBoxFuncMap(&naPools, workerState, s)

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
		boxFuncMap["getServiceStateId"] = gopcp.ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
			return workerState.State.StateId, nil
		})

		return gopcp.GetSandbox(boxFuncMap)
	}, workerStartConf)

	if !stdWorkerConfig.Nonblocking {
		utils.RunForever()
	}

	return naPools
}
