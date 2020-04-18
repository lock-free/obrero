package mids

import (
	"fmt"
	"github.com/lock-free/gopcp"
)

func RecoverMid(prefix string, fn gopcp.GeneralFun) gopcp.GeneralFun {
	return func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (ret interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%s, recover from %v", prefix, r)
			}
		}()
		ret, err = fn(args, attachment, pcpServer)
		return
	}
}
