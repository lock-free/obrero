// Package dlb: dynamic load balancer
package dlb

import (
	"math/rand"
	"sync"
)

type Worker struct {
	Id     string
	Group  string
	idx    int
	Handle interface{}
}

// Load balancer for workers
type WorkerLB struct {
	// {group:id:worker}
	ActiveWorkerMap map[string]map[string]*Worker
	// {group:[]worker}
	ActiveWorkers map[string][]*Worker
	lock          sync.Mutex
}

func GetWorkerLB() *WorkerLB {
	return &WorkerLB{make(map[string]map[string]*Worker), make(map[string][]*Worker), sync.Mutex{}}
}

func (wlb *WorkerLB) AddWorker(worker Worker) {
	wlb.lock.Lock()
	defer wlb.lock.Unlock()

	// push to the tail
	worker.idx = len(wlb.ActiveWorkers[worker.Group])
	wlb.ActiveWorkers[worker.Group] = append(wlb.ActiveWorkers[worker.Group], &worker)

	// update map
	if _, ok := wlb.ActiveWorkerMap[worker.Group]; !ok {
		wlb.ActiveWorkerMap[worker.Group] = make(map[string]*Worker)
	}
	wlb.ActiveWorkerMap[worker.Group][worker.Id] = &worker
}

func (wlb *WorkerLB) RemoveWorker(worker Worker) bool {
	wlb.lock.Lock()
	defer wlb.lock.Unlock()

	if _, ok := wlb.ActiveWorkerMap[worker.Group][worker.Id]; !ok {
		return false
	}

	// use last worker overrides current worker
	lastWorker := wlb.ActiveWorkers[worker.Group][len(wlb.ActiveWorkers[worker.Group])-1]
	lastWorker.idx = worker.idx
	wlb.ActiveWorkers[worker.Group][lastWorker.idx] = lastWorker
	// remove last one
	wlb.ActiveWorkers[worker.Group] = wlb.ActiveWorkers[worker.Group][:len(wlb.ActiveWorkers[worker.Group])-1]

	// update map
	delete(wlb.ActiveWorkerMap[worker.Group], worker.Id)
	if len(wlb.ActiveWorkerMap[worker.Group]) == 0 {
		delete(wlb.ActiveWorkerMap, worker.Group)
		delete(wlb.ActiveWorkers, worker.Group)
	}
	return true
}

func (wlb *WorkerLB) PickUpWorkerRandom(group string) (*Worker, bool) {
	wlb.lock.Lock()
	defer wlb.lock.Unlock()

	n := len(wlb.ActiveWorkers[group])
	if n <= 0 {
		return nil, false
	}

	idx := rand.Intn(n)
	return wlb.ActiveWorkers[group][idx], true
}

func (wlb *WorkerLB) PickUpWorkerById(group string, workerId string) (*Worker, bool) {
	wlb.lock.Lock()
	defer wlb.lock.Unlock()

	if _, ok := wlb.ActiveWorkerMap[group]; !ok {
		return nil, false
	}
	worker, ok := wlb.ActiveWorkerMap[group][workerId]
	return worker, ok
}

// TODO support round-robin
