package utils

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	var (
		a int
		b string
		c bool
	)
	ParseArgs([]interface{}{12, "123", true}, []interface{}{&a, &b, &c}, "")
	AssertEqual(t, a, 12, "")
	AssertEqual(t, b, "123", "")
	AssertEqual(t, c, true, "")
}

func TestParseArgs2(t *testing.T) {
	a := make(map[string]interface{})
	ParseArgs([]interface{}{map[string]interface{}{
		"key":  "v1",
		"key2": 12,
	}}, []interface{}{&a}, "")
	AssertEqual(t, a["key"], "v1", "")
	AssertEqual(t, a["key2"], float64(12), "")
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
