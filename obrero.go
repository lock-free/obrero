package obrero

import (
	"errors"
	"github.com/idata-shopee/gopcp_rpc"
	"github.com/idata-shopee/gopool"
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
	PoolSize      int
	Duration      time.Duration
	RetryDuration time.Duration
}

// Define a worker by passing `generateSandbox` function
func StartBlockWorker(generateSandbox gopcp_rpc.GenerateSandbox, workerStartConf WorkerStartConf) {
	StartWorker(generateSandbox, workerStartConf)
	// blocking forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func StartWorker(generateSandbox gopcp_rpc.GenerateSandbox, workerStartConf WorkerStartConf) []*gopool.Pool {
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

	return pools
}

func MustEnvOption(envName string) string {
	if v := os.Getenv(envName); v == "" {
		panic("missing env " + envName + " which must exists.")
	} else {
		return v
	}
}
