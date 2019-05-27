package obrero

import (
	"errors"
	"github.com/idata-shopee/gopcp_rpc"
	"log"
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

func MaintainConnectionWithNA(NAHost string, NAPort int, generateSandbox gopcp_rpc.GenerateSandbox) {
	log.Printf("try to connect to NA %s:%d\n", NAHost, NAPort)
	_, err := gopcp_rpc.GetPCPRPCClient(NAHost, NAPort, generateSandbox, func(err error) {
		log.Printf("connection to NA %s:%d broken, error is %v\n", NAHost, NAPort, err)
		time.Sleep(2 * time.Second)
		MaintainConnectionWithNA(NAHost, NAPort, generateSandbox)
	})

	if err != nil {
		log.Printf("fail to connect to NA %s:%d\n", NAHost, NAPort)
		time.Sleep(2 * time.Second)
		MaintainConnectionWithNA(NAHost, NAPort, generateSandbox)
	} else {
		log.Printf("connected to NA %s:%d\n", NAHost, NAPort)
	}
}

func StartBlockWorker(generateSandbox gopcp_rpc.GenerateSandbox) {
	StartWorker(generateSandbox)
	// blocking forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func StartWorker(generateSandbox gopcp_rpc.GenerateSandbox) {
	nas, err := ParseNAs(MustEnvOption("NAs"))
	if err != nil {
		panic(err)
	}

	if len(nas) < 0 {
		panic(errors.New("missing NAs config"))
	}

	for _, na := range nas {
		MaintainConnectionWithNA(na.Host, na.Port, generateSandbox)
	}
}

func MustEnvOption(envName string) string {
	if v := os.Getenv(envName); v == "" {
		panic("missing env " + envName + " which must exists.")
	} else {
		return v
	}
}
