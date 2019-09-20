package dlb

import (
	"github.com/lock-free/obrero/utils"
	"testing"
)

func TestBase(t *testing.T) {
	wlb := GetWorkerLB()

	w1 := Worker{Id: "0", Group: "a"}
	w2 := Worker{Id: "1", Group: "a"}
	wlb.AddWorker(w1)
	wlb.AddWorker(w2)

	for i := 0; i < 100; i++ {
		_, ok := wlb.PickUpWorkerRandom("a")
		utils.AssertEqual(t, ok, true, "")
	}

	for i := 0; i < 100; i++ {
		_, ok := wlb.PickUpWorkerRandom("b")
		utils.AssertEqual(t, ok, false, "")
	}

	wlb.RemoveWorker(w1)
	for i := 0; i < 100; i++ {
		_, ok := wlb.PickUpWorkerRandom("a")
		utils.AssertEqual(t, ok, true, "")
	}
	wlb.RemoveWorker(w2)
	for i := 0; i < 100; i++ {
		_, ok := wlb.PickUpWorkerRandom("a")
		utils.AssertEqual(t, ok, false, "")
	}
}

func TestRemoveFalsy(t *testing.T) {
	wlb := GetWorkerLB()

	w1 := Worker{Id: "0", Group: "a"}
	w2 := Worker{Id: "1", Group: "a"}
	utils.AssertEqual(t, wlb.RemoveWorker(w1), false, "")
	wlb.AddWorker(w2)
	utils.AssertEqual(t, wlb.RemoveWorker(w2), true, "")
	wlb.AddWorker(w1)
	utils.AssertEqual(t, wlb.RemoveWorker(w1), true, "")
}
