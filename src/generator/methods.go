package generator

import (
	"encoding/json"

	"github.com/lukasjarosch/go-docx"
)

func New(templateFileName string, json_data string) FileGenerator {
	activeTemplate, err := docx.Open(templateFileName)
	if err != nil {
		panic(err)
	}
	fileGenerator := FileGenerator{TempateFilename: templateFileName}
	fileGenerator.activeTemplate = activeTemplate
	fileGenerator.data = parseJson(json_data)
	return fileGenerator
}

func (s *FileGenerator) GenerateZip(filename string) {}

func parseJson(json_data string) ParseData {
	parseData := ParseData{}
	json.Unmarshal([]byte(json_data), &parseData)
	return parseData
}
