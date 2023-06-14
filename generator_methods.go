package yadt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"

	"golang.org/x/sync/errgroup"

	"github.com/tymbaca/go-yadt/utils"

	"github.com/lukasjarosch/go-docx"
)

var err error

// FileGenerator constructor. Takes template and json data as an io.Reader's.
func New(templateStream io.Reader, jsonStream io.Reader) (*FileGenerator, error) {
	templateBytes, err := utils.StreamToBytes(templateStream)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := utils.StreamToBytes(jsonStream)
	if err != nil {
		return nil, err
	}
	fileGenerator, err := NewFromBytes(templateBytes, jsonBytes)
	if err != nil {
		return nil, err
	}

	return fileGenerator, nil
}

// FileGenerator constructor. Takes template and json data in byte slices.
func NewFromBytes(templateBytes []byte, jsonBytes []byte) (*FileGenerator, error) {
	fileGenerator := new(FileGenerator)

	fileGenerator.templateBytes = templateBytes
	fileGenerator.data, err = parseJsonToData(jsonBytes)
	if err != nil {
		return nil, err
	}

	return fileGenerator, nil
}

// FileGenerator constructor. Takes template and json data filenames and read them.
func NewFromFiles(templateFilename string, jsonFilename string) (*FileGenerator, error) {

	templateBytes, err := os.ReadFile(templateFilename)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := os.ReadFile(jsonFilename)
	if err != nil {
		return nil, err
	}

	fileGenerator, err := NewFromBytes(templateBytes, jsonBytes)
	if err != nil {
		return nil, err
	}

	return fileGenerator, nil
}

// Generates docx files and packs them into zip file at the specified path.
func (s *FileGenerator) GenerateZip(path string) error {
	var err error
	s.tmpDirectory, err = os.MkdirTemp("", "") // TODO CHANGE MkdirTemp PARAMETERS TO ("", "") AFTER FIX
	defer os.RemoveAll(s.tmpDirectory)
	if err != nil {
		panic(errors.New("Error while creating temporary directory: " + err.Error()))
	}

	err = s.generateFiles()
	if err != nil {
		return errors.New("Generation error: " + err.Error())
	}
	err = utils.CompressFiles(s.filenames, path)
	if err != nil {
		return errors.New("Compression error: " + err.Error())
	}

	return nil
}

func (s *FileGenerator) generateFiles() error {
	s.filenames = []string{}
	errg, _ := errgroup.WithContext(context.Background())
	for i, fileData := range *s.data {
		currentFileData := fileData

		resultFilename := path.Join(s.tmpDirectory, (*s.data)[i].Filename+".docx")
		errg.Go(func() error {
			err := generateFile(currentFileData, s.templateBytes, resultFilename, s.tmpDirectory)
			return err
		})
		s.filenames = append(s.filenames, resultFilename)
	}
	return errg.Wait()
}

func generateFile(fileData fileData, templateBytes []byte, resultFilename string, tmpDirectory string) error {
	var pageFilenames []string

	for i, pageData := range fileData.Pages {
		pageFilename := path.Join(tmpDirectory, fileData.Filename+"_"+fmt.Sprint(i)+".docx")
		pageFilenames = append(pageFilenames, pageFilename)

		generatePageFile(templateBytes, pageFilename, pageData)
	}

	// Avoiding pagemerger call if it unnecessary
	if len(pageFilenames) >= 2 {
		err := mergePageFilesToFile(pageFilenames, resultFilename)
		if err != nil {
			return err
		}
	} else if len(pageFilenames) == 1 {
		err := os.Rename(pageFilenames[0], resultFilename)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("file with name '%s' does not have page data", fileData.Filename)
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
		return errors.New("there is no specified files to merge, pass 1 or more filenames")
	}
	args := append([]string{mergerSetPageBreaksOption, mergedFilename}, targetFilenames...)
	mergerCommand := exec.Command(pageMergerName, args...)
	output, err := mergerCommand.Output()
	if err != nil {
		log.Println(string(output))
		panic(err)
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
