package generator

import (
	"github.com/lukasjarosch/go-docx"
)

type FileGenerator struct {
	TempateFilename string
	activeTemplate  *docx.Document
	data            ParseData
}

type ParseData []FileData

type FileData struct {
	Filename string                `json:"filename"`
	Pages    []docx.PlaceholderMap `json:"pages"`
}

func NewFileGenerator(templateFileName string, json_data string) FileGenerator {
	activeTemplate, err := docx.Open(templateFileName)
	if err != nil {
		panic(err)
	}
	fileGenerator := FileGenerator{TempateFilename: templateFileName}
	fileGenerator.activeTemplate = activeTemplate
	fileGenerator.data = parseJson(json_data)
	return fileGenerator
}
