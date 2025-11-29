package main

import (
	"context"
	"fmt"
	"nas-rclone/load_env"
	"nas-rclone/state"
	"os/exec"

	"github.com/robfig/cron/v3"
)

func main() {
	fmt.Println("[INFO] start application")
	initializationEnv, err := load_env.NewInitializationEnv()
	if err != nil {
		panic(err)
	}
	if len(initializationEnv.StorageConfigurations) < 1 {
		panic("[init][ERROR] no configuration")
	}

	workerPool := cron.New()
	workerBlocker := state.NewWorkerBlocker(len(initializationEnv.StorageConfigurations))
	for i, storageConfig := range initializationEnv.StorageConfigurations {
		if _, err := workerPool.AddFunc(storageConfig.Cron, func() {
			workerState := workerBlocker.WorkerStates[i]
			if workerState.IsRunning() {
				fmt.Printf("[INFO][worker: %d] skipped\n", i+1)
				return
			}

			workerState.SetIsRunning(true)
			fmt.Printf("[INFO][worker: %d] start sync job\n", i+1)
			ctx := context.Background()

			cmd := exec.CommandContext(
				ctx,
				"./rclone", "sync",
				storageConfig.LocalRootPath, // source
				fmt.Sprintf(
					"%s:%s",
					storageConfig.Alias,
					storageConfig.RemoteRootPath,
				), // destination
				"-P",
				"--create-empty-src-dirs",
				"--transfers", fmt.Sprintf("%d", storageConfig.Transfers),
				"--checksum",
			)

			if err := cmd.Run(); err != nil {
				panic(err.Error())
			}

			workerState.SetIsRunning(false)
			fmt.Printf("[INFO][worker: %d] end sync job\n", i+1)
		}); err != nil {
			panic("[init][ERROR] failed to create cron job")
		}
	}

	workerPool.Run()
	fmt.Println("[INFO] stop application")
}
