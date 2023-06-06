package generator

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"templater/utils"

	"github.com/lukasjarosch/go-docx"
)

const MERGER_PROGRAM_NAME string = "pagemerger"
const MERGER_PROGRAM_SET_PAGEBREAKS_OPTION string = "-b"

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
	s.GenerateFiles()
	utils.CompressFiles(s.filenames, filename)
}

func (s *FileGenerator) GenerateFiles() {
	s.filenames = []string{}
	for _, fileData := range s.data {
		filename := fileData.generateFile(s.activeTemplate)
		s.filenames = append(s.filenames, filename)
	}
}

func (s *FileData) generateFile(template *docx.Document) string {
	var pageFilenames []string
	for i, pageData := range s.Pages {
		pageFilename := s.Filename + "_" + fmt.Sprint(i) + ".docx"
		generatePageFile(template, pageFilename, pageData)
	}
	resultFilename := s.Filename + ".docx"
	mergePageFilesToFile(pageFilenames, resultFilename)
	return resultFilename
}

func generatePageFile(template *docx.Document, pageFilename string, pageData docx.PlaceholderMap) {
	template.ReplaceAll(pageData)
	template.WriteToFile(pageFilename)
}

func mergePageFilesToFile(targetFilenames []string, mergedFilename string) error {
	// args := append([]string{MERGER_PROGRAM_SET_PAGEBREAKS_OPTION, mergedFilename}, targetFilenames)
	args := append([]string{MERGER_PROGRAM_SET_PAGEBREAKS_OPTION, mergedFilename}, targetFilenames...)
	mergerCommand := exec.Command(MERGER_PROGRAM_NAME, args...)
	err := mergerCommand.Run()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func parseJson(json_data string) ParseData {
	parseData := ParseData{}
	json.Unmarshal([]byte(json_data), &parseData)
	return parseData
}
