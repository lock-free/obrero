package mids

import (
	"github.com/lock-free/gopcp"
	"github.com/lock-free/gopcp_rpc"
	"github.com/lock-free/obrero/utils"
	"time"
)

type GetWorkerHandlerFun func(string) (*gopcp_rpc.PCPConnectionHandler, error)

type ProxyMid struct {
	GetWorkerHandler GetWorkerHandlerFun
}

func GetProxyMid(getWorkerHandler GetWorkerHandlerFun) *ProxyMid {
	return &ProxyMid{getWorkerHandler}
}

// lazy sandbox
// (xxx, serviceType, exp, timeout)
func (this *ProxyMid) Proxy(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
	// parse params
	var (
		serviceType string
		exp         interface{}
		timeout     int
	)

	err := utils.ParseArgs(args, []interface{}{&serviceType, &exp, &timeout}, "wrong signature, expect (xxx, serviceType: string, exp, timeout: int)")
	exp = args[1]

	if err != nil {
		return nil, err
	}

	// pick worker handle
	handle, err := this.GetWorkerHandler(serviceType)
	if err != nil {
		return nil, err
	}

	// convert exp to json string
	bs, err := gopcp.JSONMarshal(gopcp.ParseAstToJsonObject(exp))
	if err != nil {
		return nil, err
	}

	return handle.CallRemote(string(bs), time.Duration(timeout)*time.Second)
}
