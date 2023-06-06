package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"templater/utils"

	"github.com/lukasjarosch/go-docx"
)

const TMP_DIRECTORY string = "tmp/"

const MERGER_PROGRAM_NAME string = "pagemerger"
const MERGER_PROGRAM_SET_PAGEBREAKS_OPTION string = "-b"

func New(templateFilename string, json_data string) (*FileGenerator, error) {
	fileGenerator := new(FileGenerator)
	_, err := os.Stat(templateFilename)
	if err != nil {
		return nil, err
	}
	fileGenerator.TempateFilename = templateFilename
	fileGenerator.data, err = parseJson(json_data)
	return fileGenerator, nil
}

func (s *FileGenerator) GenerateZip(filename string) error {

	err := s.GenerateFiles()
	if err != nil {
		return err
	}
	err = utils.CompressFiles(s.filenames, filename)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileGenerator) GenerateFiles() error {
	s.filenames = []string{}
	for _, fileData := range *s.data {
		filename, err := fileData.generateFileAndReturnFilename(s.TempateFilename)
		if err != nil {
			return err
		}
		s.filenames = append(s.filenames, filename)
	}
	return nil
}

func (s *FileData) generateFileAndReturnFilename(templateFilename string) (string, error) {
	var pageFilenames []string
	for i, pageData := range s.Pages {
		pageFilename := TMP_DIRECTORY + s.Filename + "_" + fmt.Sprint(i) + ".docx"
		pageFilenames = append(pageFilenames, pageFilename)
		err := generatePageFile(templateFilename, pageFilename, pageData)
		if err != nil {
			return "", err
		}
	}
	resultFilename := TMP_DIRECTORY + s.Filename + ".docx"
	err := mergePageFilesToFile(pageFilenames, resultFilename)
	if err != nil {
		return "", err
	}
	return resultFilename, nil
}

func generatePageFile(templateFilename string, outputFilename string, pageData docx.PlaceholderMap) error {
	template, err := docx.Open(templateFilename)
	if err != nil {
		return err
	}

	err = template.ReplaceAll(pageData)
	if err != nil {
		return err
	}

	err = template.WriteToFile(outputFilename)
	if err != nil {
		return err
	}
	return nil
}

func mergePageFilesToFile(targetFilenames []string, mergedFilename string) error {
	// args := append([]string{MERGER_PROGRAM_SET_PAGEBREAKS_OPTION, mergedFilename}, targetFilenames)
	if len(targetFilenames) < 1 {
		return errors.New("There is no specified files to merge. Pass 1 or more filenames.")
	}
	args := append([]string{MERGER_PROGRAM_SET_PAGEBREAKS_OPTION, mergedFilename}, targetFilenames...)
	mergerCommand := exec.Command(MERGER_PROGRAM_NAME, args...)
	err := mergerCommand.Run()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func parseJson(json_data string) (*ParseData, error) {
	parseData := new(ParseData)
	err := json.Unmarshal([]byte(json_data), &parseData)
	if err != nil {
		return nil, err
	}
	return parseData, nil
}
