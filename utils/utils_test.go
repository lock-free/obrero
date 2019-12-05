package utils

import (
	"reflect"
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

func TestParseArgMap(t *testing.T) {
	var a int
	var b string
	ParseArgMap(map[string]interface{}{"a": 123, "b": "hello"}, map[string]interface{}{"a": &a, "b": &b}, "")
	AssertEqual(t, a, 123, "")
	AssertEqual(t, b, "hello", "")
}

func TestPick(t *testing.T) {
	AssertEqual(t,
		reflect.DeepEqual(Pick(map[string]interface{}{"a": 1, "b": 2, "c": 3}, []string{}), map[string]interface{}{}),
		true,
		"")
	AssertEqual(t,
		reflect.DeepEqual(Pick(map[string]interface{}{"a": 1, "b": 2, "c": 3}, []string{"a", "b"}), map[string]interface{}{"a": 1, "b": 2}),
		true,
		"")
	AssertEqual(t,
		reflect.DeepEqual(Pick(map[string]interface{}{"a": 1, "b": 2, "c": 3}, []string{"a", "b", "e"}), map[string]interface{}{"a": 1, "b": 2}),
		true,
		"")
}

func TestAssign(t *testing.T) {
	AssertEqual(t,
		reflect.DeepEqual(Assign(map[string]interface{}{"a": 1, "b": 2}, map[string]interface{}{"c": 3}), map[string]interface{}{"a": 1, "b": 2, "c": 3}),
		true,
		"")
}
