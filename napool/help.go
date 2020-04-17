package napool

import (
	"github.com/lock-free/gopcp"
	"time"
)

var pcpClient = gopcp.PcpClient{}

func (naPools *NAPools) SimpleCall(serviceType string, methodName string, args ...interface{}) (interface{}, error) {
	return naPools.CallProxy(serviceType, pcpClient.Call(methodName, args), 120*time.Second)
}
