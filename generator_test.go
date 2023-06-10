package yadt

import (
	"fmt"
	"testing"
)

var (
	waybillFilename   = "tests/waybill.json"
	templateFilename  = "tests/test_template.docx"
	outputZipFilename = "tests/output.zip"
)

func TestGenerateZip(t *testing.T) {
	fileGenerator, err := NewFromFiles(templateFilename, waybillFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	fileCount := len(*fileGenerator.data)
	pageCount := 0
	for _, fileData := range *fileGenerator.data {
		pageCount += len(fileData.Pages)
	}

	t.Log(fmt.Sprintf("Starting generation for %d files, %d pages total...", fileCount, pageCount))
	t.Log(fileGenerator.tmpDirectory)
	err = fileGenerator.GenerateZip(outputZipFilename)
	if err != nil {
		t.Errorf(err.Error())
	}
}
