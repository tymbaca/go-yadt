package yadt

import "os"

const mergerEnvKey string = "PAGEMERGER_NAME"
const mergerDefaultProgramName string = "pagemerger"
const mergerProgramSetPageBreaksOption string = "-b"

var mergerProgramName string

func init() {
	envMergerName := os.Getenv(mergerEnvKey)
	if envMergerName == "" {
		mergerProgramName = mergerDefaultProgramName
	} else {
		mergerProgramName = envMergerName
	}
}
