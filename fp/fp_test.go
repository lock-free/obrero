package fp

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
			[]interface{}{1, 2, 3},
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

func TestMap(t *testing.T) {
	ans, err := Map([]interface{}{1, 2, 3}, func(v interface{}) (interface{}, error) {
		return v.(int) + 1, nil
	})

	utils.AssertEqual(t, err, nil, "")
	utils.AssertEqual(t, ans, []interface{}{2, 3, 4}, "")
}

func TestFilter(t *testing.T) {
	ans, err := Filter([]interface{}{1, 2, 3}, func(v interface{}) (bool, error) {
		return v.(int)+1 > 2, nil
	})

	utils.AssertEqual(t, err, nil, "")
	utils.AssertEqual(t, ans, []interface{}{2, 3}, "")
}
