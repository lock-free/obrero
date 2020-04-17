package dta

import (
	"github.com/lock-free/obrero/utils"
	"testing"
)

func TestGet(t *testing.T) {
	var getData = []struct {
		Des    string
		Input  []interface{}
		Path   string
		Result interface{}
		Err    error
	}{
		{
			"list-0",
			List(1, 2, 3),
			"0",
			1,
			nil,
		},
	}

	for _, d := range getData {
		t.Run(d.Des, func(t *testing.T) {
			v, err := Get(d.Input, d.Path)
			utils.AssertEqual(t, err, d.Err, "")
			utils.AssertEqual(t, v, d.Result, "")
		})
	}
}
