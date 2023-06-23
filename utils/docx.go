package utils

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
)

var (
	ErrPlaceholdersNotFound = errors.New("placeholders not found in document")
	ErrDelimitersNotPassed  = errors.New("delimiters did not passed")
	xmlTextTag              = "t"
)

/*
FindPlaceholders finds all placeholders inside of given docx file (as byte slice).
Takes delimiter regex pattern. For example if placeholder is "{{my_key}}", then pattern would be "{{.*}}"
If given pattern is empty string then func would use default pattern "{.*}".

Returns found
placeholders as a string slice. Returns error:
- If there were no placeholders
- If given byte slice is not a valid docx file
*/
func FindPlaceholders(templateBytes []byte, leftDelimiter string, rightDelimiter string) ([]string, error) {
	if leftDelimiter == "" || rightDelimiter == "" {
		return nil, ErrDelimitersNotPassed
	}

	documentXml, err := getDocumentXml(templateBytes)
	if err != nil {
		return nil, err
	}
	var placeholders []string

	xmlDecoder := xml.NewDecoder(documentXml)

	delimiterRegexPattern := fmt.Sprintf("%s(.*)%s", leftDelimiter, rightDelimiter)
	regex, err := regexp.Compile(delimiterRegexPattern)
	if err != nil {
		return nil, err
	}
	var textElements []xml.StartElement
	for {
		token, _ := xmlDecoder.Token()
		if token == nil {
			break // End of file
		}
		textElements = appendTextElements(token, textElements)
	}

	placeholders = findPlaceholders(textElements, xmlDecoder, regex)

	if len(placeholders) > 0 {
		return placeholders, nil
	} else {
		return nil, ErrPlaceholdersNotFound
	}
}

func appendTextElements(element interface{}, textElements []xml.StartElement) []xml.StartElement {
	if startElement, ok := element.(xml.StartElement); ok {
		if len(startElement.Attr) > 0 {
			for _, attr := range startElement.Attr {
				textElements = appendTextElements(attr, textElements)
			}
		}
		if startElement.Name.Local == xmlTextTag {
			textElements = append(textElements, startElement)
		}
	}
	return textElements
}

func findPlaceholders(textElements []xml.StartElement, xmlDecoder *xml.Decoder, regex *regexp.Regexp) []string {
	var placeholders []string
	for _, element := range textElements {
		var innerText string
		xmlDecoder.DecodeElement(&innerText, &element)
		if match := regex.MatchString(innerText); match {
			placeholderKey := regex.FindStringSubmatch(innerText)[1]
			placeholders = append(placeholders, placeholderKey)
		}
	}
	return placeholders
}

var documentXmlPathInZip = "word/document.xml"

// Gets `word/document.xml` file inside of given `docx` file (it's basically `zip`). Returns error
// if file is not a valid `docx`
func getDocumentXml(templateBytes []byte) (io.Reader, error) {
	templateReader := bytes.NewReader(templateBytes)
	// Opening docx as a zip
	zipReader, err := zip.NewReader(templateReader, int64(len(templateBytes)))
	if err != nil {
		return nil, err
	}

	documentXml, err := zipReader.Open(documentXmlPathInZip)
	if err != nil {
		return nil, err
	}
	return documentXml, nil
}
