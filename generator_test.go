package yadt

import (
	"fmt"
	"os"
	"testing"
)

var (
	waybillFilename   = "tests/waybill.json"
	templateFilename  = "tests/test_template.docx"
	outputZipFilename = "tests/output.zip"
)

func TestGenerateZip(t *testing.T) {
	body, err := os.ReadFile(waybillFilename)
	fileGenerator, err := New(templateFilename, body)
	if err != nil {
		t.Errorf(err.Error())
	}

	fileCount := len(*fileGenerator.Data)
	pageCount := 0
	for _, fileData := range *fileGenerator.Data {
		pageCount += len(fileData.Pages)
	}

	t.Log(fmt.Sprintf("Starting generation for %d files, %d pages total...", fileCount, pageCount))
	t.Log(fileGenerator.tmpDirectory)
	err = fileGenerator.GenerateZip(outputZipFilename)
	if err != nil {
		t.Errorf(err.Error())
	}
}
