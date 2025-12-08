package load_env

import (
	"encoding/json"
	"fmt"
	"nas-rclone/common"
	"os"
)

type SyncType string

const (
	SyncTypeCopy SyncType = "copy"
	SyncTypeSync SyncType = "sync"
)

var (
	allSyncTypes = common.NewImmutableSet[SyncType](SyncTypeCopy, SyncTypeSync)
)

type workerConfiguration struct {
	Alias           string   `json:"alias"`
	WorkerName      string   `json:"workerName"`
	SyncType        SyncType `json:"syncType"`
	DestinationPath string   `json:"destinationPath"`
	SourcePath      string   `json:"sourcePath"`
	Cron            string   `json:"cron"`
	Transfers       int      `json:"transfers"`
}

type InitializationEnv struct {
	WorkerConfigurations []workerConfiguration `json:"workerConfigurations"`
}

func NewInitializationEnv() (*InitializationEnv, error) {
	configBytes, err := os.ReadFile(".config/.config.json")
	if err != nil {
		return nil, fmt.Errorf("[init][ERROR] failed to load config file")
	}

	var config InitializationEnv
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("[init][ERROR] failed to parse configuration")
	}

	// validation
	for _, config := range config.WorkerConfigurations {
		if !allSyncTypes.Contains(config.SyncType) {
			return nil, fmt.Errorf("[init][ERROR] invalid sync type")
		}
	}

	return &config, nil
}
