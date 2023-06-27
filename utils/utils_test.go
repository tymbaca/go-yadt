package utils

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	imageFilename                = "tests/image.png"
	articleFilename              = "tests/article.md"
	soundFilename                = "tests/sound.mp3"
	goodTemplate                 = "tests/template.docx"
	templateWithEmptyPlaceholder = "tests/empty_template.docx"
	templateWithoutPlaceholders  = "tests/template_without_placeholders.docx"
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

func TestFindPlaceholdersInDocx(t *testing.T) {
	templateBytes, err := os.ReadFile(goodTemplate)
	if err != nil {
		panic(err)
	}
	placeholders, err := FindPlaceholders(templateBytes, "{", "}")
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(placeholders, []string{"organisation", "address"}) {
		t.Fail()
	}
}

func TestFindEmptyPlaceholder(t *testing.T) {
	templateBytes, err := os.ReadFile(templateWithEmptyPlaceholder)
	if err != nil {
		panic(err)
	}
	_, err = FindPlaceholders(templateBytes, "{", "}")
	if !errors.Is(err, ErrTemplatePlaceholdersNotFound) {
		t.Fail()
	}
}

func TestPlaceholdersNotFound(t *testing.T) {
	templateBytes, err := os.ReadFile(templateWithoutPlaceholders)
	if err != nil {
		panic(err)
	}
	_, err = FindPlaceholders(templateBytes, "{", "}")
	if !errors.Is(err, ErrTemplatePlaceholdersNotFound) {
		t.Fail()
	}
}

func TestUniqueStringSlice(t *testing.T) {
	slice := []string{"hello", "hello", "you", "hello", "mama"}
	unique_slice := UniqueStringSlice(slice)
	reflect.DeepEqual(unique_slice, []string{"mama", "you", "hello"})
}

func TestCompareSlicesAsSets(t *testing.T) {
	slice1 := []string{"hello", "world"}
	slice2 := []string{"world", "world", "hello", "world", "hello"}
	if CompareSlicesAsSets[string](slice1, slice2) != true {
		t.Fail()
	}
}
