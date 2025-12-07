package load_env

import (
	"encoding/json"
	"fmt"
	"os"
)

type storageConfiguration struct {
	Alias           string `json:"alias"`
	WorkerName      string `json:"workerName"`
	DestinationPath string `json:"destinationPath"`
	SourcePath      string `json:"sourcePath"`
	Cron            string `json:"cron"`
	Transfers       int    `json:"transfers"`
}

type InitializationEnv struct {
	StorageConfigurations []storageConfiguration `json:"storageConfigurations"`
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

	return &config, nil
}
