package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func MustEnvOption(envName string) string {
	if v := os.Getenv(envName); v == "" {
		panic("missing env " + envName + " which must exists.")
	} else {
		return v
	}
}

func MustEnvIntOption(envName string) int {
	intv, err := strconv.Atoi(MustEnvOption(envName))
	if err != nil {
		panic("Env PORT must be a number.")
	}
	return intv
}

func ExistsFile(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ReadJson(filePath string, f interface{}) error {
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(source), f)
}

func WriteJson(filePath string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, bs, 0644)
}

// parse args and assign values to pointers
func ParseArgs(args []interface{}, ps []interface{}, errMsg string) error {
	if len(args) < len(ps) {
		return fmt.Errorf("missing some args, args=%v, %s", args, errMsg)
	}

	for i, p := range ps {
		err := ParseArg(args[i], p)
		if err != nil {
			return fmt.Errorf("fail to parse arg at %d, args=%v, %s", i, args, errMsg)
		}
	}
	return nil
}

// @param argMap arg as a map
// @param pm point map
func ParseArgMap(argMap map[string]interface{}, pm map[string]interface{}, errMsg string) error {
	for key, p := range pm {
		v, ok := argMap[key]
		if !ok {
			return fmt.Errorf("fail to parse arg at %s, argMap=%v, %s", key, argMap, errMsg)
		}
		err := ParseArg(v, p)
		if err != nil {
			return fmt.Errorf("fail to parse arg at %s, argMap=%v, %s", key, argMap, errMsg)
		}
	}
	return nil
}

func ParseArg(arg interface{}, pointer interface{}) error {
	bs, err := json.Marshal(arg)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, pointer)
}

func RunForever() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func AssertEqual(t *testing.T, expect interface{}, actual interface{}, message string) {
	if expect == actual {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("expect %v !=  actual %v", expect, actual)
	}
	t.Fatal(message)
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Assign(target map[string]interface{}, source map[string]interface{}) map[string]interface{} {
	ans := make(map[string]interface{})

	for k, v := range target {
		ans[k] = v
	}

	for k, v := range source {
		ans[k] = v
	}

	return ans
}

func Pick(m map[string]interface{}, fields []string) map[string]interface{} {
	ans := make(map[string]interface{})

	for _, field := range fields {
		if v, ok := m[field]; ok {
			ans[field] = v
		}
	}

	return ans
}
