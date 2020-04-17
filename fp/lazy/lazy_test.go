package lazy

import (
	"github.com/lock-free/obrero/utils"
	"testing"
)

func TestMap(t *testing.T) {
	ans, err := Eva(
		Map([]interface{}{1, 2, 3}, func(v interface{}) (interface{}, error) {
			return v.(int) + 1, nil
		}),
	)

	utils.AssertEqual(t, err, nil, "")
	utils.AssertEqual(t, ans, []interface{}{2, 3, 4}, "")
}

func TestFilter(t *testing.T) {
	ans, err := Eva(
		Filter(
			Map([]interface{}{1, 2, 3}, func(v interface{}) (interface{}, error) {
				return v.(int) + 1, nil
			}),

			func(v interface{}) (bool, error) {
				return v.(int) > 2, nil
			},
		),
	)

	utils.AssertEqual(t, err, nil, "")
	utils.AssertEqual(t, ans, []interface{}{3, 4}, "")
}

func TestGet(t *testing.T) {
	ans, err := Eva(
		Get([]interface{}{1, 2, 3}, "1"),
	)

	utils.AssertEqual(t, err, nil, "")
	utils.AssertEqual(t, ans, 2, "")
}

func TestFunctor(t *testing.T) {
	ans, err := Eva(
		Functor(Data(10, nil), func(v interface{}) (interface{}, error) {
			return v.(int) + 10, nil
		}),
	)

	utils.AssertEqual(t, err, nil, "")
	utils.AssertEqual(t, ans, 20, "")
}
