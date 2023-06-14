package yadt

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	waybillFilename    = "tests/waybill.json"
	badWaybillFilename = "tests/bad_waybill.json"
	templateFilename   = "tests/test_template.docx"
	outputZipFilename  = "tests/output.zip"
)

func TestNew(t *testing.T) {
	templateStream, _ := os.Open(templateFilename)
	jsonStream, _ := os.Open(waybillFilename)

	templateBytesReference, _ := os.ReadFile(templateFilename)
	jsonBytesReference, _ := os.ReadFile(waybillFilename)

	fileGenerator, err := New(templateStream, jsonStream)
	if err != nil {
		t.Errorf(err.Error())
	}

	var parseDataReference *parseData
	json.Unmarshal(jsonBytesReference, &parseDataReference)

	assert.Equal(t, fileGenerator.templateBytes, templateBytesReference)
	assert.Equal(t, fileGenerator.data, parseDataReference)
}

func TestNewFromBytes(t *testing.T) {
	templateBytesReference, _ := os.ReadFile(templateFilename)
	jsonBytesReference, _ := os.ReadFile(waybillFilename)

	fileGenerator, err := NewFromBytes(templateBytesReference, jsonBytesReference)
	if err != nil {
		t.Errorf(err.Error())
	}

	var parseDataReference *parseData
	json.Unmarshal(jsonBytesReference, &parseDataReference)

	assert.Equal(t, fileGenerator.templateBytes, templateBytesReference)
	assert.Equal(t, fileGenerator.data, parseDataReference)
}

func TestNewFromFiles(t *testing.T) {
	fileGenerator, err := NewFromFiles(templateFilename, waybillFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	templateBytesReference, _ := os.ReadFile(templateFilename)

	var parseDataReference *parseData
	jsonBytes, _ := os.ReadFile(waybillFilename)
	json.Unmarshal(jsonBytes, &parseDataReference)

	assert.Equal(t, fileGenerator.templateBytes, templateBytesReference)
	assert.Equal(t, fileGenerator.data, parseDataReference)
}

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

	t.Logf("Starting generation for %d files, %d pages total...", fileCount, pageCount)
	t.Log(fileGenerator.tmpDirectory)

	// Run test
	err = fileGenerator.GenerateZip(outputZipFilename)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestBadGenerateZip(t *testing.T) {
	fileGenerator, err := NewFromFiles(templateFilename, badWaybillFilename)
	if err != nil {
		t.Errorf(err.Error())
	}

	fileCount := len(*fileGenerator.data)
	pageCount := 0
	for _, fileData := range *fileGenerator.data {
		pageCount += len(fileData.Pages)
	}

	t.Logf("Starting generation for %d files, %d pages total...", fileCount, pageCount)
	t.Log(fileGenerator.tmpDirectory)

	// Run test
	err = fileGenerator.GenerateZip(outputZipFilename)
	// This is too smell
	if err != nil && strings.Contains(err.Error(), "does not have page data") {
		// PASS
		t.Log(err.Error())
	} else {
		t.Fail()
	}
}
