package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

func MustEnvOption(envName string) string {
	if v := os.Getenv(envName); v == "" {
		panic("missing env " + envName + " which must exists.")
	} else {
		return v
	}
}

func ReadJson(filePath string, f interface{}) error {
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(source), f)
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
