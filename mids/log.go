package mids

import (
	"fmt"
	"github.com/lock-free/goklog"
	"github.com/lock-free/gopcp"
	"time"
)

var klog = goklog.GetInstance()

func LogMid(logPrefix string, fn gopcp.GeneralFun) gopcp.GeneralFun {
	return func(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (ret interface{}, err error) {
		t1 := time.Now().UnixNano()

		klog.LogNormal(fmt.Sprintf("%s-access", logPrefix), fmt.Sprintf("args=%v", args))
		ret, err = fn(args, attachment, pcpServer)

		if err != nil {
			klog.LogError(fmt.Sprintf("%s-error", logPrefix), err)
		}

		t2 := time.Now().UnixNano()
		klog.LogNormal(fmt.Sprintf("%s-done", logPrefix), fmt.Sprintf("args=%v, time=%dms", args, (t2-t1)/int64(time.Millisecond)))
		return
	}
}
