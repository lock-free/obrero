package mids

import (
	"fmt"
	"github.com/lock-free/gopcp"
	"github.com/lock-free/gopcp_rpc"
	"github.com/lock-free/gopcp_stream"
	"github.com/lock-free/obrero/utils"
	"time"
)

type GetWorkerHandlerFun func(string, string) (*gopcp_rpc.PCPConnectionHandler, error)
type GetCommandFun func(gopcp.FunNode, string, int, interface{}, *gopcp.PcpServer) (string, error)

type ProxyMid struct {
	GetWorkerHandler GetWorkerHandlerFun
	GetCommand       GetCommandFun
}

func DefaultGetCommand(exp gopcp.FunNode, serviceType string, timeout int, attachment interface{}, pcpServer *gopcp.PcpServer) (string, error) {
	// convert exp to json string
	bs, err := gopcp.JSONMarshal(gopcp.ParseAstToJsonObject(exp))
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func GetProxyMid(getWorkerHandler GetWorkerHandlerFun, getCommand GetCommandFun) *ProxyMid {
	return &ProxyMid{getWorkerHandler, getCommand}
}

// lazy sandbox
// (xxx, serviceType, exp, timeout)
func (this *ProxyMid) Proxy(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
	// parse params
	var (
		serviceType string
		exp         gopcp.FunNode
		timeout     int
	)

	err := utils.ParseArgs(args, []interface{}{&serviceType, &exp, &timeout}, "wrong signature, expect (xxx, serviceType: string, exp, timeout: int)")

	if err != nil {
		return nil, err
	}

	fmt.Printf("args %v, serviceType %s, exp %v, timeout %d\n", args, serviceType, exp, timeout)

	// pick worker handle
	handle, err := this.GetWorkerHandler(serviceType, "")
	if err != nil {
		return nil, err
	}

	// translate exp to command
	cmd, err := this.GetCommand(exp, serviceType, timeout, attachment, pcpServer)
	if err != nil {
		return nil, err
	}

	// call real service
	return handle.CallRemote(cmd, time.Duration(timeout)*time.Second)
}

// lazy sandbox
// (xxx, serviceType, workerId, exp, timeout)
func (this *ProxyMid) ProxyById(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
	// parse params
	var (
		serviceType string
		workerId    string
		exp         gopcp.FunNode
		timeout     int
	)

	err := utils.ParseArgs(args, []interface{}{&serviceType, &workerId, &exp, &timeout}, "wrong signature, expect (xxx, serviceType: string, workerId: string, exp, timeout: int)")

	if err != nil {
		return nil, err
	}

	// pick worker handle
	handle, err := this.GetWorkerHandler(serviceType, workerId)
	if err != nil {
		return nil, err
	}

	// translate exp to command
	cmd, err := this.GetCommand(exp, serviceType, timeout, attachment, pcpServer)
	if err != nil {
		return nil, err
	}

	// call real service
	return handle.CallRemote(cmd, time.Duration(timeout)*time.Second)
}

// LazyStreamApi
// (xxxx, serviceType, exp, timeout)
func (this *ProxyMid) ProxyStream(streamProducer gopcp_stream.StreamProducer, args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
	// parse params
	var (
		serviceType string
		exp         gopcp.FunNode
		timeout     int
	)

	err := utils.ParseArgs(args, []interface{}{&serviceType, &exp, &timeout}, "wrong signature, expect (proxy, serviceType: string, exp, timeout: int)")

	if err != nil {
		return nil, err
	}

	var timeoutD = time.Duration(timeout) * time.Second

	jsonObj := gopcp.ParseAstToJsonObject(exp)

	switch arr := jsonObj.(type) {
	case []interface{}:
		// choose worker
		handle, err := this.GetWorkerHandler(serviceType, "")
		if err != nil {
			return nil, err
		}

		// pipe stream
		sparams, err := handle.StreamClient.ParamsToStreamParams(append(arr[1:], func(t int, d interface{}) {
			// write response of stream back to client
			switch t {
			case gopcp_stream.STREAM_DATA:
				streamProducer.SendData(d, timeoutD)
			case gopcp_stream.STREAM_END:
				streamProducer.SendEnd(timeoutD)
			default:
				errMsg, ok := d.(string)
				if !ok {
					streamProducer.SendError(fmt.Sprintf("errored at stream, and responsed error message is not string. d=%v", d), timeoutD)
				} else {
					streamProducer.SendError(errMsg, timeoutD)
				}
			}
		}))

		if err != nil {
			return nil, err
		}

		// send a stream request to service
		return handle.Call(gopcp.CallResult{append([]interface{}{arr[0]}, sparams...)}, timeoutD)
	default:
		return nil, fmt.Errorf("Expect array, but got %v, args=%v", jsonObj, args)
	}
}
