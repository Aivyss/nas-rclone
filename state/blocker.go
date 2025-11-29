package state

import "sync"

type workerState struct {
	isRunning bool
	sync.Mutex
}

func (s *workerState) IsRunning() bool {
	return s.isRunning
}

func (s *workerState) SetIsRunning(isRunning bool) {
	s.Lock()
	s.isRunning = isRunning
	s.Unlock()
}

type WorkerBlocker struct {
	WorkerStates []*workerState
}

func NewWorkerBlocker(size int) *WorkerBlocker {
	states := make([]*workerState, 0, size)
	for i := 0; i < size; i++ {
		states = append(states, new(workerState))
	}

	return &WorkerBlocker{
		WorkerStates: states,
	}
}
