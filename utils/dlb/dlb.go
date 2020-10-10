// Package dlb: dynamic load balancer
package dlb

import (
	"math/rand"
	"sync"
)

type Worker struct {
	Id     string      // identity of worker
	Group  string      // for each worker, they belongs to one group
	Info   string      // some extra information
	idx    int         `json:"-"`
	Handle interface{} `json:"-"`
}

// Load balancer for workers
type WorkerLB struct {
	// {group:id:worker}
	activeWorkerMap map[string]map[string]*Worker
	// {group:[]worker}
	activeWorkers map[string][]*Worker
	lock          sync.Mutex
}

func GetWorkerLB() *WorkerLB {
	return &WorkerLB{make(map[string]map[string]*Worker), make(map[string][]*Worker), sync.Mutex{}}
}

func (wlb *WorkerLB) AddWorker(worker *Worker) {
	wlb.lock.Lock()
	defer wlb.lock.Unlock()

	// push to the tail
	worker.idx = len(wlb.activeWorkers[worker.Group])
	wlb.activeWorkers[worker.Group] = append(wlb.activeWorkers[worker.Group], worker)

	// update map
	if _, ok := wlb.activeWorkerMap[worker.Group]; !ok {
		wlb.activeWorkerMap[worker.Group] = make(map[string]*Worker)
	}
	wlb.activeWorkerMap[worker.Group][worker.Id] = worker
}

func (wlb *WorkerLB) RemoveWorker(worker *Worker) bool {
	wlb.lock.Lock()
	defer wlb.lock.Unlock()

	if _, ok := wlb.activeWorkerMap[worker.Group][worker.Id]; !ok {
		return false
	}

	// use last worker overrides current worker
	lastWorker := wlb.activeWorkers[worker.Group][len(wlb.activeWorkers[worker.Group])-1]
	lastWorker.idx = worker.idx
	wlb.activeWorkers[worker.Group][lastWorker.idx] = lastWorker
	// remove last one
	wlb.activeWorkers[worker.Group] = wlb.activeWorkers[worker.Group][:len(wlb.activeWorkers[worker.Group])-1]

	// update map
	delete(wlb.activeWorkerMap[worker.Group], worker.Id)
	if len(wlb.activeWorkerMap[worker.Group]) == 0 {
		delete(wlb.activeWorkerMap, worker.Group)
		delete(wlb.activeWorkers, worker.Group)
	}
	return true
}

func (wlb *WorkerLB) PickUpWorkerRandom(group string) (*Worker, bool) {
	wlb.lock.Lock()
	defer wlb.lock.Unlock()

	n := len(wlb.activeWorkers[group])
	if n <= 0 {
		return nil, false
	}

	idx := rand.Intn(n)
	return wlb.activeWorkers[group][idx], true
}

func (wlb *WorkerLB) PickUpWorkerById(group string, workerId string) (*Worker, bool) {
	wlb.lock.Lock()
	defer wlb.lock.Unlock()

	if _, ok := wlb.activeWorkerMap[group]; !ok {
		return nil, false
	}
	worker, ok := wlb.activeWorkerMap[group][workerId]
	return worker, ok
}

func (wlb *WorkerLB) GetActiveWorkerMap() map[string]map[string]*Worker {
	return wlb.activeWorkerMap
}

// TODO support round-robin
