package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	imageFilename   = "tests/image.png"
	articleFilename = "tests/article.md"
	soundFilename   = "tests/sound.mp3"
)

func TestStreamToBytes(t *testing.T) {
	textBytesReference, _ := os.ReadFile(articleFilename)

	textStream, _ := os.Open(articleFilename)
	textBytes, _ := StreamToBytes(textStream)
	assert.Equal(t, textBytes, textBytesReference)
}

func TestCompressFiles(t *testing.T) {
	filesToCompress := []string{imageFilename, articleFilename, soundFilename}
	outputFilename := "tests/output.zip"

	err := CompressFiles(filesToCompress, outputFilename)
	if err != nil {
		t.Errorf(err.Error())
	}
	if _, err = os.Stat(outputFilename); err != nil {
		t.Errorf("Error: There is no output zip!")
	}
}
