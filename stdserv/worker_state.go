package stdserv

import (
	"github.com/lock-free/obrero/utils"
	"github.com/satori/go.uuid"
	"sync"
)

// Worker can have state:
// (1) unique state_id for each worker instance
// (2) load state when started
// (3) update state when running worker

type WorkerState struct {
	StateFilePath string
	State         State
	lock          *sync.Mutex
}

type State struct {
	StateId string      `json: "stateId"`
	Data    interface{} `json:"data"`
}

// read state from state file
func GetWorkerState(stateFilePath string) (*WorkerState, error) {
	// if no state file, create a new one
	if !utils.ExistsFile(stateFilePath) {
		// create initial state
		err := createInitialState(stateFilePath)
		if err != nil {
			return nil, err
		}
	}

	// if there is a state file, read from it
	var state State
	err := utils.ReadJson(stateFilePath, &state)
	if err != nil {
		return nil, err
	}

	var lock sync.Mutex
	workerState := &WorkerState{
		StateFilePath: stateFilePath,
		State:         state,
		lock:          &lock,
	}

	return workerState, nil
}

// flush current state to file
// TODO lock
func (ws *WorkerState) UpdateState() error {
	ws.lock.Lock()
	defer ws.lock.Unlock()

	return utils.WriteJson(ws.StateFilePath, ws.State)
}

func createInitialState(stateFilePath string) error {
	stateId := uuid.NewV4().String()
	state := State{StateId: stateId}
	return utils.WriteJson(stateFilePath, state)
}
