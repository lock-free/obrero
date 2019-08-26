package obrero

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
	Pools             []*gopool.Pool
	GetClientMaxRetry int
}

func (naPools *NAPools) CallProxy(serviceType string, list gopcp.CallResult, timeout time.Duration) (interface{}, error) {
	client, err := naPools.GetRandomItem()

	if err != nil {
		return nil, err
	}

	return client.Call(client.PcpClient.Call("proxy", serviceType, client.PcpClient.Call("'", list), timeout.Seconds()), timeout)
}

func (naPools *NAPools) CallProxyById(serviceId string, serviceType string, list gopcp.CallResult, timeout time.Duration) (interface{}, error) {
	client, err := naPools.GetRandomItem()

	if err != nil {
		return nil, err
	}

	return client.Call(client.PcpClient.Call("proxyById", serviceId, serviceType, client.PcpClient.Call("'", list), timeout.Seconds()), timeout)
}

func (naPools *NAPools) CallProxyStream(serviceType string, list gopcp.CallResult, streamCallback gopcp_stream.StreamCallbackFunc, timeout time.Duration) (interface{}, error) {
	client, err := naPools.GetRandomItem()

	if err != nil {
		return nil, err
	}

	sexp, err := client.StreamClient.StreamCall("proxyStream", serviceType, client.PcpClient.Call("'", list), timeout.Seconds(), streamCallback)
	if err != nil {
		return nil, err
	}

	return client.Call(*sexp, timeout)
}

// pick up a random item
func (naPools *NAPools) GetRandomItem() (*gopcp_rpc.PCPConnectionHandler, error) {
	return naPools.getRandomItem(0, naPools.GetClientMaxRetry)
}

func (naPools *NAPools) hashNA(key []byte) (*gopcp_rpc.PCPConnectionHandler, error) {
	index := getHash(key) % len(naPools.Pools)

	item, err := naPools.Pools[index].Get()
	if err != nil {
		return nil, err
	}

	client, ok := item.(*gopcp_rpc.PCPConnectionHandler)
	if !ok {
		return nil, errors.New("unexpected error at hashNA")
	} else {
		return client, nil
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
func (naPools *NAPools) getRandomItem(tryCount int, maxCount int) (*gopcp_rpc.PCPConnectionHandler, error) {
	if tryCount > maxCount {
		return nil, errors.New("fail to get a connection from NA pools, tried 3 times")
	}

	// pick up a random na pool.
	index := rand.Intn(len(naPools.Pools))

	item, err := naPools.Pools[index].Get()

	if err != nil {
		return naPools.getRandomItem(tryCount+1, maxCount)
	} else {
		client, ok := item.(*gopcp_rpc.PCPConnectionHandler)
		if !ok {
			// TODO sleep a while before retry
			return naPools.getRandomItem(tryCount+1, maxCount)
		} else {
			return client, nil
		}
	}
}
