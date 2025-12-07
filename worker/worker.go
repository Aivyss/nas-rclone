package worker

import (
	"fmt"
	"nas-rclone/common"
	"sync"
)

var (
	IsRunningWorkerErr = fmt.Errorf("running worker")
)

type Worker interface {
	SyncRun(func() error) error
	IsRunning() bool
}

type Workers []Worker

/*
GetProgressPercent
return -> 0 ~ 100
*/
func (ws Workers) GetProgressPercent() int {
	totalWorkersSize := len(ws)
	doneWorkersSize := len(common.FilterSlice(ws, func(e Worker) bool {
		return !e.IsRunning()
	}))

	return (100 * doneWorkersSize) / totalWorkersSize
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

func (w *worker) IsRunning() bool {
	return w.isRunning
}

func NewWorker() Worker {
	return new(worker)
}
