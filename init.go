package yadt

import (
	"os"
)

const (
	mergerEnvKey              string = "PAGEMERGER_NAME"
	mergerDefaultName         string = "pagemerger"
	mergerSetPageBreaksOption string = "-b"
)

func getPageMergerPath() string {
	envMergerName := os.Getenv(mergerEnvKey)
	if envMergerName == "" {
		return mergerDefaultName
	}

	return envMergerName
}
