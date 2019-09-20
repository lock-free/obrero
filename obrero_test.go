package obrero

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

func TestParseNAs(t *testing.T) {
	v, _ := ParseNAs("127.0.0.1:8001")
	assertEqual(t, len(v), 1, "")
	assertEqual(t, v[0], NA{"127.0.0.1", 8001}, "")
}

func TestParseNAs2(t *testing.T) {
	v, _ := ParseNAs("127.0.0.1:8001;120.130.140.2:9087")
	assertEqual(t, len(v), 2, "")
	assertEqual(t, v[0], NA{"127.0.0.1", 8001}, "")
	assertEqual(t, v[1], NA{"120.130.140.2", 9087}, "")
}
