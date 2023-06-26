package yadt

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/lukasjarosch/go-docx"
	"github.com/stretchr/testify/assert"
)

var (
	// Data
	goodWaybill                = "tests/waybill.json"
	waybillWithEmptyPageData   = "tests/bad_waybill_with_empty_pagedata.json"
	waybillWithDifferentFields = "tests/bad_waybill_different_fields.json"
	waybillIncompatible        = "tests/bad_waybill_not_compatible.json"
	templateFilename           = "tests/template.docx"

	// Templates
	emptyTemplate                 = "tests/empty_template.docx"
	templateWithLeadingWhitespace = "tests/template_leading_whitespace.docx"
	templateWithTailingWhitespace = "tests/template_tailing_whitespace.docx"
	templateWithBothWhitespace    = "tests/template_both_whitespace.docx"
	brokenTemplate                = "tests/broken_template.docx"

	// Output
	outputZipFilename = "tests/output.zip"
	outputDocx        = "tests/output.docx"
)

func TestNew(t *testing.T) {

	templateStream, err := os.Open(templateFilename)
	if err != nil {
		panic(err)
	}
	jsonStream, err := os.Open(goodWaybill)
	if err != nil {
		panic(err)
	}

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
	templateBytesReference, err := os.ReadFile(templateFilename)
	if err != nil {
		panic(err)
	}
	jsonBytesReference, err := os.ReadFile(goodWaybill)
	if err != nil {
		panic(err)
	}

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

	// just for log
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

func TestEmptyPageValidate(t *testing.T) {
	_, err := NewFromFiles(templateFilename, waybillWithEmptyPageData)
	if errors.Is(err, ErrFileDataWithoutPages) {
		// Ok
	} else {
		t.Fatal()
	}
}

func TestWhitespaceValidate(t *testing.T) {
	_, err := NewFromFiles(templateWithLeadingWhitespace, goodWaybill)
	if !errors.Is(err, ErrPlaceholderWithWhitespaces) {
		t.Fatal()
	}

	_, err = NewFromFiles(templateWithTailingWhitespace, goodWaybill)
	if !errors.Is(err, ErrPlaceholderWithWhitespaces) {
		t.Fatal()
	}

	_, err = NewFromFiles(templateWithBothWhitespace, goodWaybill)
	if !errors.Is(err, ErrPlaceholderWithWhitespaces) {
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
	_, err := NewFromFiles(brokenTemplate, waybillIncompatible)
	if err == nil {
		t.Fail()
	}
}

func TestEmptyTemplate(t *testing.T) {
	_, err := NewFromFiles(emptyTemplate, goodWaybill)
	if errors.Is(err, ErrTemplatePlaceholdersNotFound) {
		//ok
	} else {
		t.Error(err)
	}
}

func TestPageFileGenerate(t *testing.T) {
	templateBytes, err := os.ReadFile(templateFilename)
	if err != nil {
		panic(err)
	}
	pageData := docx.PlaceholderMap{
		"organisation": "Baskins Bread",
		"address":      "Frunze 42",
	}
	//////////////////
	err = generatePageFile(templateBytes, outputDocx, pageData)
	if err != nil {
		t.Error(err)
	}
}
