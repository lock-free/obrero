package lazy

import (
	"errors"
	"github.com/lock-free/obrero/fp/operation"
)

type LazyValue interface {
	Eva() (interface{}, error)
}

func Eva(c interface{}) (interface{}, error) {
	switch v := c.(type) {
	case LazyValue:
		return v.Eva()
	default:
		return v, nil
	}
}

type DataLazyValue struct {
	data interface{}
	err  error
}

func (dlv DataLazyValue) Eva() (interface{}, error) {
	return dlv.data, dlv.err
}

func Data(c interface{}, err error) LazyValue {
	return DataLazyValue{
		data: c,
		err:  err,
	}
}

type FunctionLazyValue struct {
	eval func() (interface{}, error)
}

func (t FunctionLazyValue) Eva() (interface{}, error) {
	return t.eval()
}

type MapItem func(interface{}) (interface{}, error)

func Map(v interface{}, mapItem MapItem) LazyValue {
	return FunctionLazyValue{
		eval: func() (interface{}, error) {
			list, err := Eva(v)
			if err != nil {
				return nil, err
			}
			switch items := list.(type) {
			case []interface{}:
				var ans []interface{}
				for _, v := range items {
					n, err := mapItem(v)
					if err != nil {
						return nil, err
					}
					ans = append(ans, n)
				}
				return ans, nil
			default:
				return nil, errors.New("Expect []interface type")
			}
		},
	}
}

func Get(v interface{}, path string) LazyValue {
	return FunctionLazyValue{
		eval: func() (interface{}, error) {
			item, err := Eva(v)
			if err != nil {
				return nil, err
			}

			return operation.Get(item, path)
		},
	}
}

type Predicate func(interface{}) (bool, error)

func Filter(v LazyValue, predicate Predicate) LazyValue {
	return FunctionLazyValue{
		eval: func() (interface{}, error) {
			list, err := v.Eva()
			if err != nil {
				return nil, err
			}
			switch items := list.(type) {
			case []interface{}:
				var ans []interface{}
				for _, v := range items {
					// TODO recover from panic
					pass, err := predicate(v)
					if err != nil {
						return nil, err
					}
					if pass {
						ans = append(ans, v)
					}
				}
				return ans, nil
			default:
				return nil, errors.New("Expect []interface type")
			}
		},
	}
}
