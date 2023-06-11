package yadt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/tymbaca/go-yadt/utils"

	"github.com/lukasjarosch/go-docx"
)

const mergerProgramName string = "pagemerger"
const mergerProgramSetPageBreaksOption string = "-b"

var err error

func New(templateStream io.Reader, jsonStream io.Reader) (*FileGenerator, error) {
	fileGenerator := new(FileGenerator)

	templateBytes, err := utils.StreamToBytes(templateStream)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := utils.StreamToBytes(jsonStream)
	if err != nil {
		return nil, err
	}

	fileGenerator.templateBytes = templateBytes
	fileGenerator.data, err = parseJsonToData(jsonBytes)
	if err != nil {
		return nil, err
	}
	return fileGenerator, nil
}

func NewFromBytes(templateBytes []byte, jsonBytes []byte) (*FileGenerator, error) {
	fileGenerator := new(FileGenerator)

	fileGenerator.templateBytes = templateBytes
	fileGenerator.data, err = parseJsonToData(jsonBytes)
	if err != nil {
		return nil, err
	}

	return fileGenerator, nil
}

func NewFromFiles(templateFilename string, jsonFilename string) (*FileGenerator, error) {
	fileGenerator := new(FileGenerator)

	// Check if files exists
	_, err := os.Stat(templateFilename)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(jsonFilename)
	if err != nil {
		return nil, err
	}

	fileGenerator.templateBytes, err = os.ReadFile(templateFilename)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := os.ReadFile(jsonFilename)
	if err != nil {
		return nil, err
	}
	fileGenerator.data, err = parseJsonToData(jsonBytes)
	if err != nil {
		return nil, err
	}

	return fileGenerator, nil
}

func (s *FileGenerator) GenerateZip(filename string) error {
	var err error
	s.tmpDirectory, err = os.MkdirTemp("", "")
	if err != nil {
		panic(errors.New("Error while creating temporary directory: " + err.Error()))
	}
	defer os.RemoveAll(s.tmpDirectory)

	err = s.generateFiles()
	if err != nil {
		return errors.New("Generation error: " + err.Error())
	}
	err = utils.CompressFiles(s.filenames, filename)
	if err != nil {
		return errors.New("Compression error: " + err.Error())
	}

	return nil
}

func (s *FileGenerator) generateFiles() error {
	s.filenames = []string{}
	var wg sync.WaitGroup
	for i, fileData := range *s.data {
		resultFilename := path.Join(s.tmpDirectory, (*s.data)[i].Filename+".docx")
		wg.Add(1)
		go func() {
			fileData.generateFile(s.templateBytes, resultFilename, s.tmpDirectory)
			defer wg.Done()
		}()
		s.filenames = append(s.filenames, resultFilename)
	}
	wg.Wait()
	return nil
}

func (s *fileData) generateFile(templateBytes []byte, resultFilename string, tmpDirectory string) error {
	var pageFilenames []string
	var wg sync.WaitGroup
	for i, pageData := range s.Pages {
		pageFilename := tmpDirectory + s.Filename + "_" + fmt.Sprint(i) + ".docx"
		pageFilenames = append(pageFilenames, pageFilename)

		wg.Add(1)
		go func() {
			generatePageFile(templateBytes, pageFilename, pageData)
			defer wg.Done()
		}()
	}
	wg.Wait()

	err := mergePageFilesToFile(pageFilenames, resultFilename)
	if err != nil {
		return err
	}
	return nil
}

func generatePageFile(templateBytes []byte, outputFilename string, pageData docx.PlaceholderMap) error {
	template, err := docx.OpenBytes(templateBytes)
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
	args := append([]string{mergerProgramSetPageBreaksOption, mergedFilename}, targetFilenames...)
	mergerCommand := exec.Command(mergerProgramName, args...)
	err := mergerCommand.Run()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func parseJsonToData(json_data []byte) (*parseData, error) {
	parseData := new(parseData)
	err := json.Unmarshal(json_data, &parseData)
	if err != nil {
		return nil, err
	}
	return parseData, nil
}
