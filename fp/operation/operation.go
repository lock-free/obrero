package operation

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
	var cur = value

	parts := strings.Split(jsonPath, ".")

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
						return nil, errors.New("missing value for path: " + jsonPath + ". Out of range. Array length is " + strconv.Itoa(len(nextObjectParent)) + ".")
					}
					nextObject := nextObjectParent[num]
					cur = nextObject
					continue
				}
			}

			// otherwise regarding as map
			nextObjectParent, ok := cur.(map[string]interface{})
			if !ok {
				return nil, errors.New("Can not go deeper for this jsonPath: " + jsonPath + ". Type of current object is " + fmt.Sprintf("%v", reflect.TypeOf(cur)))
			}

			nextObject, ok := nextObjectParent[part]

			if !ok {
				return nil, errors.New("missing value for path: " + jsonPath)
			} else {
				cur = nextObject
			}
		}
	}

	return cur, nil
}

func List(args ...interface{}) []interface{} {
	return args
}
