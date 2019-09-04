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

func (naPools *NAPools) CallProxy(serviceType string,
	exp gopcp.CallResult,
	timeout time.Duration) (interface{}, error) {
	client, err := naPools.GetRandomItem()

	if err != nil {
		return nil, err
	}

	return naPools.CallNAProxy(client, serviceType, exp, timeout)
}

func (naPools *NAPools) CallNAProxy(client *gopcp_rpc.PCPConnectionHandler,
	serviceType string,
	exp gopcp.CallResult,
	timeout time.Duration) (interface{}, error) {
	return client.Call(client.PcpClient.Call(
		"proxy",
		serviceType,
		client.PcpClient.Call("'", exp),
		timeout.Seconds()), timeout)
}

func (naPools *NAPools) CallProxyById(serviceId string,
	serviceType string,
	exp gopcp.CallResult,
	timeout time.Duration) (interface{}, error) {
	client, err := naPools.GetRandomItem()

	if err != nil {
		return nil, err
	}

	return naPools.CallNAProxyById(client, serviceId, serviceType, exp, timeout)
}

func (naPools *NAPools) CallNAProxyById(client *gopcp_rpc.PCPConnectionHandler,
	serviceId string,
	serviceType string,
	exp gopcp.CallResult,
	timeout time.Duration) (interface{}, error) {
	return client.Call(client.PcpClient.Call(
		"proxyById",
		serviceId,
		serviceType,
		client.PcpClient.Call("'", exp),
		timeout.Seconds()), timeout)
}

func (naPools *NAPools) CallProxyStream(serviceType string,
	exp gopcp.CallResult,
	streamCallback gopcp_stream.StreamCallbackFunc,
	timeout time.Duration) (interface{}, error) {
	client, err := naPools.GetRandomItem()

	if err != nil {
		return nil, err
	}

	return naPools.CallNAProxyStream(client, serviceType, exp, streamCallback, timeout)
}

func (naPools *NAPools) CallNAProxyStream(client *gopcp_rpc.PCPConnectionHandler,
	serviceType string,
	exp gopcp.CallResult,
	streamCallback gopcp_stream.StreamCallbackFunc,
	timeout time.Duration) (interface{}, error) {
	sexp, err := client.StreamClient.StreamCall("proxyStream",
		serviceType,
		client.PcpClient.Call("'", exp),
		timeout.Seconds(),
		streamCallback)

	if err != nil {
		return nil, err
	}

	return client.Call(*sexp, timeout)
}

// pick up a random item
func (naPools *NAPools) GetRandomItem() (*gopcp_rpc.PCPConnectionHandler, error) {
	return naPools.getRandomItem(0, naPools.GetClientMaxRetry)
}

// pick up NA by hash key
func (naPools *NAPools) HashNA(key string) (*gopcp_rpc.PCPConnectionHandler, error) {
	index := getHash([]byte(key)) % len(naPools.Pools)

	item, err := naPools.Pools[index].Get()
	if err != nil {
		return nil, err
	}

	client, ok := item.(*gopcp_rpc.PCPConnectionHandler)
	if !ok {
		return nil, errors.New("unexpected error at HashNA")
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
