package utils

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var (
	ErrTemplatePlaceholdersNotFound = errors.New("placeholders not found in template")
	ErrDelimitersNotPassed          = errors.New("delimiters did not passed")
	ErrBadTemplate                  = errors.New("can't open template file, it's broken")
	ErrPlaceholderWithWhitespaces   = errors.New("some placeholders template has leading or tailing whitespace")

	documentXmlPathInZip = "word/document.xml"
	xmlTextTag           = "t"
)

/*
FindPlaceholders finds all placeholders inside of given docx file (as byte slice).
Takes delimiter regex pattern. For example if placeholder is "{{my_key}}", then pattern would be "{{.*}}"
If given pattern is empty string then func would use default pattern "{.*}".

Returns found
placeholders as a string slice. Returns error:
- If there were no placeholders
- If given byte slice is not a valid docx file
- If placeholders have leading or tailing whitespace inside of placeholder (like `{ key}`, `{key }` or `{ key }`)
*/
func FindPlaceholders(templateBytes []byte, leftDelimiter string, rightDelimiter string) ([]string, error) {
	if leftDelimiter == "" || rightDelimiter == "" {
		return nil, ErrDelimitersNotPassed
	}
	delimiterRegexPattern := fmt.Sprintf("%s(.+?)%s", leftDelimiter, rightDelimiter)

	documentXmlReader, err := getDocumentXmlReader(templateBytes)
	if err != nil {
		return nil, err
	}
	documentText := getAllXmlText(documentXmlReader)
	placeholders, err := findPlaceholders(documentText, delimiterRegexPattern)
	if err != nil {
		return nil, err
	}
	return placeholders, nil
}

func getDocumentXmlReader(templateBytes []byte) (io.Reader, error) {
	templateReader := bytes.NewReader(templateBytes)
	zipReader, err := zip.NewReader(templateReader, int64(len(templateBytes)))
	if err != nil {
		return nil, ErrBadTemplate
	}
	file, err := zipReader.Open(documentXmlPathInZip)
	if err != nil {
		return nil, ErrBadTemplate
	}
	return file, nil
}

// Gets `word/document.xml` as string from given `docx` file (it's basically `zip`). Returns error
// if file is not a valid `docx`
func getAllXmlText(reader io.Reader) string {

	var output string
	tokenizer := html.NewTokenizer(reader)
	prevToken := tokenizer.Token()
loop:
	for {
		tok := tokenizer.Next()
		switch {
		case tok == html.ErrorToken:
			break loop // End of the document,  done
		case tok == html.StartTagToken:
			prevToken = tokenizer.Token()
		case tok == html.TextToken:
			if prevToken.Data == "script" {
				continue
			}
			TxtContent := html.UnescapeString(string(tokenizer.Text()))
			if len(TxtContent) > 0 {
				output += TxtContent
			}
		}
	}
	return output
}

func findPlaceholders(text string, delimiterRegexPattern string) ([]string, error) {
	var placeholders []string
	r, err := regexp.Compile(delimiterRegexPattern)
	if err != nil {
		return nil, err
	}

	// find and process
	matchSet := r.FindAllStringSubmatch(text, -1)
	for _, submatchPair := range matchSet {
		placeholder := submatchPair[1] // submatch is second element in FindAllStringSubmatch result slice

		if placeholder != strings.TrimSpace(placeholder) {
			return nil, ErrPlaceholderWithWhitespaces
		}

		placeholders = append(placeholders, placeholder)
	}

	// return
	if len(placeholders) > 0 {
		return placeholders, nil
	} else {
		return nil, ErrTemplatePlaceholdersNotFound
	}
}
