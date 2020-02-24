package napool

import (
	"errors"
	"github.com/lock-free/gopcp"
	"github.com/lock-free/gopcp_rpc"
	"github.com/lock-free/gopcp_stream"
	"github.com/lock-free/gopool"
	"math/rand"
	"time"
)

type NAPools struct {
	// create a pool for each NA
	Pools             []*gopool.Pool
	GetClientMaxRetry int
}

// (proxy, serviveType, exp, timeout)
func (naPools *NAPools) CallProxy(serviceType string, exp gopcp.CallResult, timeout time.Duration) (interface{}, error) {
	client, err := naPools.GetRandomItem()

	if err != nil {
		return nil, err
	}

	return client.Call(client.PcpClient.Call("proxy", serviceType, exp, timeout.Seconds()), timeout)
}

// (proxyById, serviceType, workerId, exp, timeout)
func (naPools *NAPools) CallProxyById(serviceType string, workerId string, exp gopcp.CallResult, timeout time.Duration) (interface{}, error) {
	client, err := naPools.GetRandomItem()

	if err != nil {
		return nil, err
	}

	return client.Call(client.PcpClient.Call("proxyById", serviceType, workerId, exp, timeout.Seconds()), timeout)
}

func (naPools *NAPools) CallProxyStream(serviceType string, exp gopcp.CallResult, streamCallback gopcp_stream.StreamCallbackFunc, timeout time.Duration) (interface{}, error) {
	client, err := naPools.GetRandomItem()

	if err != nil {
		return nil, err
	}

	sexp, err := client.StreamClient.StreamCall("proxyStream", serviceType, exp, timeout.Seconds(), streamCallback)

	if err != nil {
		return nil, err
	}

	return client.Call(*sexp, timeout)
}

// pick up a random item
func (naPools *NAPools) GetRandomItem() (*gopcp_rpc.PCPConnectionHandler, error) {
	item, err := naPools.getRandomItem(0, naPools.GetClientMaxRetry)
	if err != nil {
		return nil, err
	}
	return itemToPcpConnectionHandler(item)
}

// pick up NA by hash key
func (naPools *NAPools) HashNA(key string) (*gopcp_rpc.PCPConnectionHandler, error) {
	index := getHash([]byte(key)) % len(naPools.Pools)

	item, err := naPools.Pools[index].Get()
	if err != nil {
		return nil, err
	}
	return itemToPcpConnectionHandler(item)
}

func itemToPcpConnectionHandler(item interface{}) (*gopcp_rpc.PCPConnectionHandler, error) {
	v, ok := item.(*gopcp_rpc.PCPConnectionHandler)
	if !ok {
		return nil, errors.New("expect type of gopcp_rpc.PCPConnectionHandler")
	} else {
		return v, nil
	}
}

func getHash(data []byte) int {
	var s uint32 = 2166136261

	// write data
	for _, c := range data {
		s *= 16777619
		s ^= uint32(c)
	}

	return int(s)
}

// TODO implement robin-round instead of random
func (naPools *NAPools) getRandomItem(tryCount int, maxCount int) (interface{}, error) {
	if tryCount > maxCount {
		return nil, errors.New("fail to get a connection from NA pools, tried 3 times")
	}

	// pick up a random na pool.
	if len(naPools.Pools) == 0 {
		return nil, errors.New("empty pools")
	}
	index := rand.Intn(len(naPools.Pools))

	item, err := naPools.Pools[index].Get()

	if err != nil {
		return naPools.getRandomItem(tryCount+1, maxCount)
	} else {
		return item, nil
	}
}
