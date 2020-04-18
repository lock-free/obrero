package dt

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// dynamic type access

// basic type: int, int32, int64, float32, float64, string, boolean
// composed type: slice, map

// get value by json path
// type explaination
func Get(value interface{}, jsonPath string) (interface{}, error) {
	parts := strings.Split(jsonPath, ".")
	return getByJsonPath(value, parts)
}

func getByJsonPath(value interface{}, parts []string) (interface{}, error) {
	var cur = value
	for _, part := range parts {
		if part != "" {
			// check if it is number
			num, err := strconv.Atoi(part)
			if err == nil {
				// try array
				nextObjectParent, ok := cur.([]interface{})
				if ok {
					// out of boundry
					if num < 0 || num > len(nextObjectParent) {
						return nil, fmt.Errorf("missing value. Out of range.")
					}
					nextObject := nextObjectParent[num]
					cur = nextObject
					continue
				}
			}

			// otherwise regarding as map
			nextObjectParent, ok := cur.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Can not go deeper. Type of current object is %v", reflect.TypeOf(cur))
			}

			nextObject, ok := nextObjectParent[part]

			if !ok {
				return nil, errors.New("missing value")
			} else {
				cur = nextObject
			}
		}
	}

	return cur, nil
}

func Set(source interface{}, jsonPath string, value interface{}) (interface{}, error) {
	parts := strings.Split(jsonPath, ".")
	if len(parts) <= 0 {
		return value, nil
	}
	obj, err := getByJsonPath(source, parts[:len(parts)-1])
	if err != nil {
		return nil, err
	}

	key := parts[len(parts)-1]

	// try array
	if num, err := strconv.Atoi(key); err != nil {
		list, ok := obj.([]interface{})
		if ok {
			if num < 0 || num > len(list) {
				return nil, fmt.Errorf("Out of range When Set. Array length is %d", len(list))
			}
			list[num] = value
		}
	}

	// try map
	m, ok := obj.(map[string]interface{})
	if !ok {
		return nil, errors.New("Expect map")
	}
	m[key] = value

	return source, nil
}

type ItemHandler func(interface{}) error

func ForEach(list interface{}, itemHandler ItemHandler) error {
	switch items := list.(type) {
	case []interface{}:
		for _, v := range items {
			err := itemHandler(v)
			if err != nil {
				return err
			}
		}
		return nil
	case map[string]interface{}:
		for _, v := range items {
			err := itemHandler(v)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New("Expect []interface type")
	}
}

type MapItem func(interface{}) (interface{}, error)

func Map(list interface{}, mapItem MapItem) ([]interface{}, error) {
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
	case map[string]interface{}:
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
}

func MapIndex(list interface{}, mapItem MapItem) ([]interface{}, error) {
	switch items := list.(type) {
	case []interface{}:
		var ans []interface{}
		for index, _ := range items {
			n, err := mapItem(index)
			if err != nil {
				return nil, err
			}
			ans = append(ans, n)
		}
		return ans, nil
	case map[string]interface{}:
		var ans []interface{}
		for index, _ := range items {
			n, err := mapItem(index)
			if err != nil {
				return nil, err
			}
			ans = append(ans, n)
		}
		return ans, nil
	default:
		return nil, errors.New("Expect []interface type")
	}
}

type Predicate func(interface{}) (bool, error)

func Filter(list interface{}, predicate Predicate) (interface{}, error) {
	switch items := list.(type) {
	case []interface{}:
		var ans []interface{}
		for _, v := range items {
			pass, err := predicate(v)
			if err != nil {
				return nil, err
			}
			if pass {
				ans = append(ans, v)
			}
		}
		return ans, nil
	case map[string]interface{}:
		var ans []interface{}
		for _, v := range items {
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
}

func FilterIndex(list interface{}, predicate Predicate) (interface{}, error) {
	switch items := list.(type) {
	case []interface{}:
		var ans []interface{}
		for index, v := range items {
			pass, err := predicate(v)
			if err != nil {
				return nil, err
			}
			if pass {
				ans = append(ans, index)
			}
		}
		return ans, nil
	case map[string]interface{}:
		var ans []interface{}
		for index, v := range items {
			pass, err := predicate(v)
			if err != nil {
				return nil, err
			}
			if pass {
				ans = append(ans, index)
			}
		}
		return ans, nil
	default:
		return nil, errors.New("Expect []interface type")
	}
}

func Falsy(v interface{}) bool {
	return v == nil || v == false || v == 0 || v == ""
}
