package main

import (
	"context"
	"errors"
	"fmt"
	"nas-rclone/load_env"
	"nas-rclone/worker"
	"os/exec"

	"github.com/robfig/cron/v3"
)

func main() {
	fmt.Println("[INFO] start application")
	initializationEnv, err := load_env.NewInitializationEnv()
	if err != nil {
		panic(err)
	}
	if len(initializationEnv.WorkerConfigurations) < 1 {
		panic("[init][ERROR] no configuration")
	}

	workerPool := cron.New()
	workers := createWorkers(len(initializationEnv.WorkerConfigurations))
	for i, workerConfig := range initializationEnv.WorkerConfigurations {
		if _, err := workerPool.AddFunc(workerConfig.Cron, func() {
			if workerRunErr := workers[i].SyncRun(func() error {
				fmt.Printf("[INFO][worker: %d][worker_name: %s] start sync job\n", i+1, workerConfig.WorkerName)
				ctx := context.Background()

				cmd := exec.CommandContext(
					ctx,
					"./rclone", string(workerConfig.SyncType),
					workerConfig.SourcePath,
					fmt.Sprintf(
						"%s:%s",
						workerConfig.Alias,
						workerConfig.DestinationPath,
					),
					"-P",
					"--create-empty-src-dirs",
					"--transfers", fmt.Sprintf("%d", workerConfig.Transfers),
					"--checksum",
					"--exclude", "._*",
					"--exclude", ".DS_Store",
				)

				return cmd.Run()
			}); workerRunErr != nil {
				if errors.Is(workerRunErr, worker.IsRunningWorkerErr) {
					fmt.Printf("[INFO][worker: %d][worker_name: %s] skipped because this worker is still running...\n", i+1, workerConfig.WorkerName)
					return
				}

				fmt.Printf("[ERROR] failed to sync job: %v\n", workerRunErr)
				return
			}

			fmt.Printf("[INFO][worker: %d][worker_name: %s] end sync job, progress: %d%% \n", i+1, workerConfig.WorkerName, workers.GetProgressPercent())
		}); err != nil {
			panic("[init][ERROR] failed to create cron job")
		}
	}

	workerPool.Run()
	fmt.Println("[INFO] stop application")
}

func createWorkers(size int) worker.Workers {
	workers := make([]worker.Worker, 0, size)
	for i := 0; i < size; i++ {
		workers = append(workers, worker.NewWorker())
	}

	return workers
}
