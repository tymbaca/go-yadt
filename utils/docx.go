package utils

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"regexp"
)

var (
	ErrPlaceholdersNotFound = errors.New("placeholders not found in document")
	xmlTextTag              = "w:t"
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
func FindPlaceholders(templateBytes []byte, delimiterRegexPattern string) ([]string, error) {
	if delimiterRegexPattern == "" {
		delimiterRegexPattern = "{.*}"
	}
	
	documentXml, err := getDocumentXml(templateBytes)
	if err != nil {
		return nil, err
	}
	var placeholders []string
	xmlDecoder := xml.NewDecoder(documentXml)
	regex := regexp.Compile(delimiterRegexPattern)
	for {
		token, _ := xmlDecoder.Token()
		if token == nil {
			break // End of file
		}

		if startElement, ok := token.(xml.StartElement); ok {
			if startElement.Name.Local == xmlTextTag {
				var innerText string
				xmlDecoder.DecodeElement(&innerText, &startElement)
				if match, _ := regex.MatchString() 
				placeholders = append(placeholders)
			}
		} else {
			return nil, errors.New("error while parsing xml")
		}
	}
	return placeholders, nil
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
