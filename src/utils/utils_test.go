package utils

import (
	"os"
	"testing"
)

var (
	filesToCompress = []string{"tests/image.png", "tests/article.md", "tests/sound.mp3"}
	outputFilename  = "tests/output.zip"
)

func TestCompressFiles(t *testing.T) {
	err := CompressFiles(filesToCompress, outputFilename)
	if err != nil {
		t.Errorf(err.Error())
	}
	if _, err = os.Stat(outputFilename); err != nil {
		t.Errorf("Error: There is no output zip!")
	}
}
