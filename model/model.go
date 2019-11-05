package model

import (
	"fmt"
	"github.com/lock-free/gopcp"
	"github.com/lock-free/obrero/napool"
	"time"
)

var pcpClient = gopcp.PcpClient{}

func CheckRel(naPools *napool.NAPools, relKey string, e1Id string, e2Id string) error {
	hasI, err := naPools.CallProxy("model_obrero", pcpClient.Call("hasRel", relKey, e1Id, e2Id), 120*time.Second)
	if err != nil {
		return err
	}

	has, ok := hasI.(bool)
	if !ok {
		return fmt.Errorf("expect bool, but get %v", hasI)
	}

	if !has {
		return fmt.Errorf("e1 %s dosen't have %s", e1Id, e2Id)
	}
	return nil
}
