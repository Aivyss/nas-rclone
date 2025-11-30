package worker

import (
	"fmt"
	"sync"
)

var (
	IsRunningWorkerErr = fmt.Errorf("running worker")
)

type Worker interface {
	SyncRun(func() error) error
}

type worker struct {
	isRunning bool
	sync.Mutex
}

func (w *worker) SyncRun(f func() error) error {
	if w.isRunning {
		return IsRunningWorkerErr
	}

	var err error
	w.Lock()
	w.isRunning = true
	err = f()
	w.isRunning = false
	w.Unlock()

	return err
}

func NewWorker() Worker {
	return new(worker)
}
