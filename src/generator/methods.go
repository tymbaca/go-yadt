package generator

import (
	"encoding/json"
	"fmt"

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

func (s *FileGenerator) GenerateZip(filename string) {
	s.generateFiles()
	s.compress(filename)
}

func (s *FileGenerator) generateFiles() {
	s.filenames = []string{}
	for _, fileData := range s.data {
		filename := fileData.generateFile(s.activeTemplate)
		s.filenames = append(s.filenames, filename)
	}
}

func (s *FileGenerator) compress(filename string) {

}

func (s *FileData) generateFile(template *docx.Document) string {
	var pageFilenames []string
	for i, pageData := range s.Pages {
		pageFilename := s.Filename + "_" + fmt.Sprint(i) + ".docx"
		generatePageFile(template, pageFilename, pageData)
	}
	resultFilename := s.Filename + ".docx"
	mergeFilesToFile(pageFilenames, resultFilename)
	return resultFilename
}

func generatePageFile(template *docx.Document, pageFilename string, pageData docx.PlaceholderMap) {
	template.ReplaceAll(pageData)
	template.WriteToFile(pageFilename)
}

func mergeFilesToFile(targetFilenames []string, mergedFilename string)

func parseJson(json_data string) ParseData {
	parseData := ParseData{}
	json.Unmarshal([]byte(json_data), &parseData)
	return parseData
}
