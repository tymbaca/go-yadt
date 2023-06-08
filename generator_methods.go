package yadt

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/tymbaca/go-yadt/utils"

	"github.com/lukasjarosch/go-docx"
)

const tmpDirectory string = "tmp/"

const MergerProgramName string = "pagemerger"
const MargerProgramSetPagebreaksOption string = "-b"

func New(templateFilename string, json_data []byte) (*FileGenerator, error) {
	fileGenerator := new(FileGenerator)
	_, err := os.Stat(templateFilename)
	if err != nil {
		return nil, err
	}
	fileGenerator.TempateFilename = templateFilename
	fileGenerator.Data, err = parseJson(json_data)

	return fileGenerator, nil
}

func (s *FileGenerator) GenerateZip(filename string) error {
	var err error
	s.tmpDirectory, err = os.MkdirTemp("", "")
	if err != nil {
		panic(errors.New("Error while creating temporary directory: " + err.Error()))
	}
	defer os.RemoveAll(s.tmpDirectory)

	err = s.GenerateFiles()
	if err != nil {
		return errors.New("Generation error: " + err.Error())
	}
	err = utils.CompressFiles(s.filenames, filename)
	if err != nil {
		return errors.New("Compression error: " + err.Error())
	}

	return nil
}

func (s *FileGenerator) GenerateFiles() error {
	s.filenames = []string{}
	var wg sync.WaitGroup
	for i, fileData := range *s.Data {
		resultFilename := path.Join(s.tmpDirectory, (*s.Data)[i].Filename+".docx")
		wg.Add(1)
		go func() {
			fileData.generateFile(s.TempateFilename, resultFilename, s.tmpDirectory)
			defer wg.Done()
		}()
		s.filenames = append(s.filenames, resultFilename)
	}
	wg.Wait()
	return nil
}

func (s *FileData) generateFile(templateFilename string, resultFilename string, tmpDirectory string) error {
	var pageFilenames []string
	var wg sync.WaitGroup
	for i, pageData := range s.Pages {
		pageFilename := tmpDirectory + s.Filename + "_" + fmt.Sprint(i) + ".docx"
		pageFilenames = append(pageFilenames, pageFilename)

		wg.Add(1)
		go func() {
			generatePageFile(templateFilename, pageFilename, pageData)
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
	args := append([]string{MargerProgramSetPagebreaksOption, mergedFilename}, targetFilenames...)
	mergerCommand := exec.Command(MergerProgramName, args...)
	err := mergerCommand.Run()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func parseJson(json_data []byte) (*ParseData, error) {
	parseData := new(ParseData)
	err := json.Unmarshal(json_data, &parseData)
	if err != nil {
		return nil, err
	}
	return parseData, nil
}
