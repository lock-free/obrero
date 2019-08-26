package obrero

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lock-free/gopcp_rpc"
	"github.com/lock-free/gopool"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// utils for worker service

// 1. parse NAs variable from env
// 2. maintain connections to NAs

// eg: NAs=127.0.0.1:8000;123.109.89.10:7000
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

// Define a worker by passing `generateSandbox` function
func StartBlockWorker(generateSandbox gopcp_rpc.GenerateSandbox, workerStartConf WorkerStartConf) {
	StartWorker(generateSandbox, workerStartConf)
	RunForever()
}

func RunForever() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

// when start a worker, will parse env variable NAs, and then
// connect to them.
func StartWorker(generateSandbox gopcp_rpc.GenerateSandbox, workerStartConf WorkerStartConf) NAPools {
	nas, err := ParseNAs(MustEnvOption("NAs"))
	if err != nil {
		panic(err)
	}

	return StartWorkerWithNAs(generateSandbox, workerStartConf, nas)
}

func StartWorkerWithNAs(generateSandbox gopcp_rpc.GenerateSandbox, workerStartConf WorkerStartConf, nas []NA) NAPools {
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

func ReadJson(filePath string, f interface{}) error {
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(source), f)
}

// parse args and assign values to pointers
func ParseArgs(args []interface{}, ps []interface{}, errMsg string) error {
	if len(args) < len(ps) {
		return fmt.Errorf("missing some args, args=%v, %s", args, errMsg)
	}

	for i, p := range ps {
		err := ParseArg(args[i], p)
		if err != nil {
			return fmt.Errorf("fail to parse arg at %d, args=%v, %s", i, args, errMsg)
		}
	}
	return nil
}

func ParseArg(arg interface{}, pointer interface{}) error {
	bs, err := json.Marshal(arg)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, pointer)
}
