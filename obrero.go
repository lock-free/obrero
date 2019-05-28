package obrero

import (
	"errors"
	"github.com/idata-shopee/gopcp"
	"github.com/idata-shopee/gopcp_rpc"
	"github.com/idata-shopee/gopool"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// utils for worker service

// 1. parse NAs variable from env
// 2. maintain connections to NAs

// NAs: 127.0.0.1:8000;123.109.89.10:7000
type NA struct {
	Host string
	Port int
}

func ParseNAs(nas string) ([]NA, error) {
	var ans []NA
	for _, naStr := range strings.Split(nas, ";") {
		parts := strings.Split(naStr, ":")
		if len(parts) <= 1 || len(parts) > 2 {
			return ans, errors.New("wrong nas str in config")
		}
		host := parts[0]
		portStr := parts[1]

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return ans, err
		}

		ans = append(ans, NA{host, port})
	}

	return ans, nil
}

type WorkerStartConf struct {
	PoolSize            int
	Duration            time.Duration
	RetryDuration       time.Duration
	NAGetClientMaxRetry int
}

type NAPools struct {
	Pools             []*gopool.Pool
	GetClientMaxRetry int
}

func (naPools *NAPools) CallProxy(serviceType string, list gopcp.CallResult, timeout time.Duration) (interface{}, error) {
	// 1. pickup a NA connection
	client, err := naPools.GetItem()
	if err != nil {
		return nil, err
	}

	code, err := client.PcpClient.ToJSON(list)
	if err != nil {
		return nil, err
	}

	// 2. call NA proxy
	return client.Call(client.PcpClient.Call("proxy", serviceType, code, timeout.Seconds()), timeout)
}

func (naPools *NAPools) GetItem() (*gopcp_rpc.PCPConnectionHandler, error) {
	return naPools.getItem(0, naPools.GetClientMaxRetry)
}

func (naPools *NAPools) getItem(tryCount int, maxCount int) (*gopcp_rpc.PCPConnectionHandler, error) {
	if tryCount > maxCount {
		return nil, errors.New("fail to get a connection from NA pools, tried 3 times")
	}

	index := rand.Intn(len(naPools.Pools))

	item, err := naPools.Pools[index].Get()

	if err != nil {
		return naPools.getItem(tryCount+1, maxCount)
	} else {
		client, ok := item.(*gopcp_rpc.PCPConnectionHandler)
		if !ok {
			return naPools.getItem(tryCount+1, maxCount)
		} else {
			return client, nil
		}
	}
}

// TODO
// func (naPools *NAPools) CallProxyStream() (interface{}, error) {
// 	// 1. pickup a
// }

// Define a worker by passing `generateSandbox` function
func StartBlockWorker(generateSandbox gopcp_rpc.GenerateSandbox, workerStartConf WorkerStartConf) {
	StartWorker(generateSandbox, workerStartConf)
	// blocking forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func StartWorker(generateSandbox gopcp_rpc.GenerateSandbox, workerStartConf WorkerStartConf) NAPools {
	nas, err := ParseNAs(MustEnvOption("NAs"))
	if err != nil {
		panic(err)
	}

	if len(nas) < 0 {
		panic(errors.New("missing NAs config"))
	}

	var pools []*gopool.Pool

	for _, na := range nas {
		pool := gopcp_rpc.GetPCPRPCPool(func() (string, int, error) {
			return na.Host, na.Port, nil
		}, generateSandbox, workerStartConf.PoolSize, workerStartConf.Duration, workerStartConf.RetryDuration)

		pools = append(pools, pool)
	}

	return NAPools{
		Pools:             pools,
		GetClientMaxRetry: workerStartConf.NAGetClientMaxRetry,
	}
}

func MustEnvOption(envName string) string {
	if v := os.Getenv(envName); v == "" {
		panic("missing env " + envName + " which must exists.")
	} else {
		return v
	}
}
