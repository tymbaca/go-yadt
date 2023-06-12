package yadt

import (
	"os"
	"os/exec"
)

const mergerEnvKey string = "PAGEMERGER_NAME"
const mergerDefaultName string = "pagemerger"
const mergerSetPageBreaksOption string = "-b"

var pageMergerName string

func init() {
	setPageMergerName()
	checkPageMergerExists()
}

func setPageMergerName() {
	envMergerName := os.Getenv(mergerEnvKey)
	if envMergerName == "" {
		pageMergerName = mergerDefaultName
	} else {
		pageMergerName = envMergerName
	}
}

func checkPageMergerExists() {
	_, err := exec.LookPath(pageMergerName)
	if err != nil {
		panic(err)
	}

}
