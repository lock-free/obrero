package cqbox

import (
	"github.com/lock-free/gopcp"
	"github.com/lock-free/obrero/utils"
	"github.com/lock-free/obrero/utils/cq"
)

type CallQueueBox struct {
	callQueueMap *cq.CallQueueMap
}

func GetCallQueueBox() *CallQueueBox {
	var callQueueMap = cq.GetCallQueueMap(cq.CALL_QUEUE_DEF_EXECUTOR)
	return &CallQueueBox{callQueueMap}
}

func (this *CallQueueBox) CallQueueBoxFn(args []interface{}, attachment interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
	var (
		key string
		exp interface{}
	)

	err := utils.ParseArgs(args, []interface{}{&key, &exp}, "wrong signature, expect (queue, key: string, exp: interface{})")

	if err != nil {
		return nil, err
	}

	// add task to queue
	return this.callQueueMap.Enqueue(key, func() (interface{}, error) {
		return pcpServer.ExecuteAst(exp, attachment)
	})
}
