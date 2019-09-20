package utils

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, expect interface{}, actual interface{}, message string) {
	if expect == actual {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("expect %v !=  actual %v", expect, actual)
	}
	t.Fatal(message)
}

func TestParseArgs(t *testing.T) {
	var (
		a int
		b string
		c bool
	)
	ParseArgs([]interface{}{12, "123", true}, []interface{}{&a, &b, &c}, "")
	assertEqual(t, a, 12, "")
	assertEqual(t, b, "123", "")
	assertEqual(t, c, true, "")
}

func TestParseArgs2(t *testing.T) {
	a := make(map[string]interface{})
	ParseArgs([]interface{}{map[string]interface{}{
		"key":  "v1",
		"key2": 12,
	}}, []interface{}{&a}, "")
	assertEqual(t, a["key"], "v1", "")
	assertEqual(t, a["key2"], float64(12), "")
}

func TestParseArgs3(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	var p Person
	ParseArgs([]interface{}{map[string]interface{}{
		"Name": "a",
		"Age":  12,
	}}, []interface{}{&p}, "")
}
