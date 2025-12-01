package yadt

import (
	"os"
)

const (
	mergerEnvKey              string = "PAGEMERGER_NAME"
	mergerDefaultName         string = "pagemerger"
	mergerSetPageBreaksOption string = "-b"
)

var pageMergerName string

func init() {
	setPageMergerName()
}

func setPageMergerName() {
	envMergerName := os.Getenv(mergerEnvKey)
	if envMergerName == "" {
		pageMergerName = mergerDefaultName
	} else {
		pageMergerName = envMergerName
	}
}
