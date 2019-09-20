package obrero

import (
	"github.com/lock-free/obrero/utils"
	"testing"
)

func TestParseNAs(t *testing.T) {
	v, _ := ParseNAs("127.0.0.1:8001")
	utils.AssertEqual(t, len(v), 1, "")
	utils.AssertEqual(t, v[0], NA{"127.0.0.1", 8001}, "")
}

func TestParseNAs2(t *testing.T) {
	v, _ := ParseNAs("127.0.0.1:8001;120.130.140.2:9087")
	utils.AssertEqual(t, len(v), 2, "")
	utils.AssertEqual(t, v[0], NA{"127.0.0.1", 8001}, "")
	utils.AssertEqual(t, v[1], NA{"120.130.140.2", 9087}, "")
}
