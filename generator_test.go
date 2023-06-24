package yadt

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	goodWaybill                = "tests/waybill.json"
	waybillWithEmptyPageData   = "tests/bad_waybill_with_empty_pagedata.json"
	waybillWithDifferentFields = "tests/bad_waybill_different_fields.json"
	waybillIncompatible        = "tests/bad_waybill_not_compatible.json"
	templateFilename           = "tests/test_template.docx"
	emptyTemplatePath          = "tests/empty_template.docx"
	outputZipFilename          = "tests/output.zip"
)

func TestNew(t *testing.T) {
	templateStream, _ := os.Open(templateFilename)
	jsonStream, _ := os.Open(goodWaybill)

	templateBytesReference, _ := os.ReadFile(templateFilename)
	jsonBytesReference, _ := os.ReadFile(goodWaybill)

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
	jsonBytesReference, _ := os.ReadFile(goodWaybill)

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
	fileGenerator, err := NewFromFiles(templateFilename, goodWaybill)
	if err != nil {
		t.Errorf(err.Error())
	}

	templateBytesReference, _ := os.ReadFile(templateFilename)

	var parseDataReference *parseData
	jsonBytes, _ := os.ReadFile(goodWaybill)
	json.Unmarshal(jsonBytes, &parseDataReference)

	assert.Equal(t, fileGenerator.templateBytes, templateBytesReference)
	assert.Equal(t, fileGenerator.data, parseDataReference)
}

func TestGenerateZip(t *testing.T) {
	fileGenerator, err := NewFromFiles(templateFilename, goodWaybill)
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

func TestEmptyPageGenerateZip(t *testing.T) {
	_, err := NewFromFiles(templateFilename, waybillWithEmptyPageData)
	if errors.Is(err, ErrFileDataWithoutPages) {
		// Ok
	} else {
		t.Fatal()
	}
}

func TestDifferentFields(t *testing.T) {
	_, err := NewFromFiles(templateFilename, waybillWithDifferentFields)
	if !errors.Is(err, ErrDataWithDifferentFields) {
		t.Fatal(err)
	}
}

func TestIncompatible(t *testing.T) {
	_, err := NewFromFiles(templateFilename, waybillIncompatible)
	if err != nil {
		// Ok
	} else {
		t.Fail()
	}
}

func TestBadTemplate(t *testing.T) {

}

func TestEmptyTemplate(t *testing.T) {
	_, err := NewFromFiles(emptyTemplatePath, goodWaybill)
	if errors.Is(err, ErrTemplatePlaceholdersNotFound) {
		//ok
	} else {
		t.Fail()
	}
}
